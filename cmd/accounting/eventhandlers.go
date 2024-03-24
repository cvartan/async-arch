package main

import (
	model "async-arch/internal/domain/accounting"
	eventmodel "async-arch/internal/domain/event"
	"async-arch/internal/lib/base"
	events "async-arch/internal/lib/event"
	"database/sql"
	"encoding/json"
	"log"
	"math/rand"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
)

func initEventHandlers() {
	eventConsumerCUD.AddConsumedEvent(eventmodel.AUTH_CUD_USER_CREATED, "1", handleUserCreatedEvent)
	eventConsumerCUD.Consume()

	eventConsumerBE.AddConsumedEvent(eventmodel.TASK_BE_TASK_CREATED, "1", handleTaskCreateEvent)
	eventConsumerBE.AddConsumedEvent(eventmodel.TASK_BE_TASK_ASSIGNED, "1", handleTaskAssignedEvent)
	eventConsumerBE.AddConsumedEvent(eventmodel.TASK_BE_TASK_COMPLETED, "1", handleTaskCompletedEvent)
	eventConsumerBE.Consume()
}

// Обработчик CUD-события добавления нового пользователя
func handleUserCreatedEvent(event *events.Event, data interface{}) {
	log.Printf("Event %s(id=\"%s\",dataid=\"%s\") received\n", event.EventType, event.EventID, event.DataID)
	// Получаем тело соообщения с событием
	body := string(data.([]byte))
	var userData eventmodel.UserEventData
	err := json.NewDecoder(strings.NewReader(body)).Decode(&userData)
	if err != nil {
		log.Fatal(err)
	}

	// Получаем данные пользователя и сохраняем их в БД
	user := &model.User{
		Uuid: userData.Uuid,
		Name: userData.Name,
	}

	repo, _ := base.App.GetDomainRepository("accounting")
	err = repo.Append(user)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Event %s(id=\"%s\") catched\n", event.EventType, event.EventID)
}

// Обработчик бизнес-события добавления новой задачи
/*
	Добавляем списание средств с пользователя после формирования цены
	TODO: когда-нибудь вынести бизнес-логику в отдельную обработку, чтобы не было дублирования операций
*/
func handleTaskCreateEvent(event *events.Event, data interface{}) {
	log.Printf("Event %s(id=\"%s\",dataid=\"%s\") received\n", event.EventType, event.EventID, event.DataID)
	// Получаем тело соообщения с событием
	body := string(data.([]byte))
	taskData := &eventmodel.PricedTaskStreamData{}
	err := json.NewDecoder(strings.NewReader(body)).Decode(taskData)
	if err != nil {
		log.Fatalf("(dataid=%s) %v\n", event.DataID, err)
	}

	task := &model.Task{
		Uuid:                taskData.Uuid,
		AssignedUserUuid:    taskData.AssignedUserUuid,
		AssignmentTaskPrice: 10 + rand.Intn(11),
		CompleteTaskPrice:   20 + rand.Intn(21),
	}

	repo, _ := base.App.GetDomainRepository("accounting")

	taskData.AssignmentTaskPrice = task.AssignmentTaskPrice
	taskData.CompletedTaskPrice = task.CompleteTaskPrice
	taskData.State = "ACTIVE"

	// Выполняем списание денег с назначенного пользователя

	trx := &model.Transaction{
		Uuid:     uuid.NewString(),
		Type:     model.DEBITING,
		Time:     time.Now(),
		UserUuid: task.AssignedUserUuid,
		TaskUuid: task.Uuid,
		Value:    task.AssignmentTaskPrice,
	}

	repo.Append(task)
	repo.Append(trx)

	trxData := &eventmodel.TransactionEventData{
		Uuid:           trx.Uuid,
		Type:           string(trx.Type),
		Time:           trx.Time,
		LinkedUserUuid: trx.UserUuid,
		LinkedTaskUuid: trx.TaskUuid,
		Value:          trx.Value,
	}

	_, err = eventProducerTaskCUD.ProduceEventData(eventmodel.ACC_CUD_TASK_PRICED, task.Uuid, reflect.TypeOf(*task).String(), taskData, "1", nil)
	if err != nil {
		log.Fatalf("event %s (dataid=%s) %v\n", eventmodel.ACC_CUD_TASK_PRICED, event.DataID, err)
	}

	_, err = eventProducerTaskCUD.ProduceEventData(eventmodel.ACC_CUD_TRX_CREATED, trx.Uuid, reflect.TypeOf(*trx).String(), trxData, "1", nil)
	if err != nil {
		log.Fatalf("event %s (dataid=%s) %v\n", eventmodel.ACC_CUD_TRX_CREATED, event.DataID, err)
	}

	_, err = eventProducerTrxBE.ProduceEventData(eventmodel.ACC_BE_DEBITING, trx.Uuid, reflect.TypeOf(*trx).String(), trxData, "1", nil)
	if err != nil {
		log.Fatalf("event %s (dataid=%s) %v\n", eventmodel.ACC_BE_DEBITING, event.DataID, err)
	}

	log.Printf("Event %s(id=\"%s\") catched\n", event.EventType, event.EventID)
}

// Обработчик бизнес-события назначения задачи пользователю
func handleTaskAssignedEvent(event *events.Event, data interface{}) {
	log.Printf("Event %s(id=\"%s\",dataid=\"%s\") received\n", event.EventType, event.EventID, event.DataID)
	// Получаем тело соообщения с событием
	body := string(data.([]byte))
	var taskData eventmodel.TaskEventData
	err := json.NewDecoder(strings.NewReader(body)).Decode(&taskData)
	if err != nil {
		log.Fatalf("(dataid=%s) %v\n", event.DataID, err)
	}

	repo, _ := base.App.GetDomainRepository("accounting")

	// Ищем задачу
	task := &model.Task{}
	err = repo.Get(task, map[string]interface{}{"uuid": taskData.Uuid})
	if err != nil {
		log.Fatalf("(dataid=%s) %v\n", event.DataID, err)
	}

	task.AssignedUserUuid = taskData.AssignedUserUuid

	trx := &model.Transaction{
		Uuid:     uuid.NewString(),
		Time:     time.Now(),
		UserUuid: task.AssignedUserUuid,
		TaskUuid: task.Uuid,
		Type:     model.DEBITING,
		Value:    task.AssignmentTaskPrice,
	}

	err = repo.Update(task)
	if err != nil {
		log.Fatalln()
	}

	err = repo.Append(trx)
	if err != nil {
		log.Fatalln(err)
	}

	trxEventData := &eventmodel.TransactionEventData{
		Uuid:           trx.Uuid,
		Type:           string(trx.Type),
		Time:           trx.Time,
		LinkedUserUuid: trx.UserUuid,
		LinkedTaskUuid: trx.TaskUuid,
		Value:          trx.Value,
	}

	_, err = eventProducerTrxCUD.ProduceEventData(eventmodel.ACC_CUD_TRX_CREATED, trx.Uuid, reflect.TypeOf(*trx).String(), trxEventData, "1", nil)
	if err != nil {
		log.Fatalln(err)
	}

	// Генерим бизнес-событие по транзакции списания
	_, err = eventProducerTrxBE.ProduceEventData(eventmodel.ACC_BE_DEBITING, trx.Uuid, reflect.TypeOf(*trx).String(), trxEventData, "1", nil)
	if err != nil {
		log.Fatalf("event %s (dataid=%s) %v\n", eventmodel.ACC_BE_DEBITING, event.DataID, err)
	}

	log.Printf("Event %s(id=\"%s\") catched\n", event.EventType, event.EventID)
}

// Обработчик бизнес-события завершения задачи пользователем
func handleTaskCompletedEvent(event *events.Event, data interface{}) {
	log.Printf("Event %s(id=\"%s\") received\n", event.EventType, event.EventID)
	// Получаем тело соообщения с событием
	body := string(data.([]byte))
	var taskData eventmodel.TaskEventData
	err := json.NewDecoder(strings.NewReader(body)).Decode(&taskData)
	if err != nil {
		log.Fatal(err)
	}

	repo, _ := base.App.GetDomainRepository("accounting")

	// Ищем задачу
	task := &model.Task{}
	err = repo.Get(task, map[string]interface{}{"uuid": taskData.Uuid})
	if err != nil {
		log.Fatalln(err)
	}

	trx := &model.Transaction{
		Uuid:     uuid.NewString(),
		Time:     time.Now(),
		UserUuid: task.AssignedUserUuid,
		TaskUuid: task.Uuid,
		Type:     model.VALUE,
		Value:    task.CompleteTaskPrice,
	}

	err = repo.Append(trx)
	if err != nil {
		log.Fatalln(err)
	}

	trxEventData := &eventmodel.TransactionEventData{
		Uuid:           trx.Uuid,
		Type:           string(trx.Type),
		Time:           trx.Time,
		LinkedUserUuid: trx.UserUuid,
		LinkedTaskUuid: trx.TaskUuid,
		Value:          trx.Value,
	}

	_, err = eventProducerTrxCUD.ProduceEventData(eventmodel.ACC_CUD_TRX_CREATED, trx.Uuid, reflect.TypeOf(*trx).String(), trxEventData, "1", nil)
	if err != nil {
		log.Fatalln(err)
	}

	// Генерим бизнес-событие по транзакции ззачисления
	_, err = eventProducerTrxBE.ProduceEventData(eventmodel.ACC_BE_VALUE, trx.Uuid, reflect.TypeOf(*trx).String(), trxEventData, "1", nil)
	if err != nil {
		log.Fatalf("event %s (dataid=%s) %v\n", eventmodel.ACC_BE_VALUE, event.DataID, err)
	}

	log.Printf("Event %s(id=\"%s\") catched\n", event.EventType, event.EventID)
}

const getUserListQueryTemplate = `
select u.id , u."uuid" ,u."name" ,u.balance from accounting."user" u
`
const getUserTransactionSumQueryTemplate = `
select 
	(select sum(t.value) from accounting."transaction" t where t."type" = 'DEBITING' and date_trunc('day',t."time")=current_date and t.user_uuid = ?) debiting_sum, 
	(select sum(t.value) from accounting."transaction" t where t."type" = 'VALUE' and date_trunc('day',t."time")=current_date and t.user_uuid = ?) value_sum 
`

// Отлавливаем событие для ребаланса и выполняем ребаланс
func handleRebalanceMessage(key, value []byte, headers map[string]interface{}) {
	repo, _ := base.App.GetDomainRepository("accounting")
	result, err := repo.RawQuery(getUserListQueryTemplate)
	if err != nil {
		log.Fatalln(err)
	}

	rows, ok := result.(*sql.Rows)
	if !ok {
		log.Fatalln("result is not rows")
	}

	var users []*model.User

	for rows.Next() {
		user := &model.User{}
		err := rows.Scan(&user.ID, &user.Uuid, &user.Name, &user.Balance)
		if err != nil {
			log.Fatalln(err)
		}

		users = append(users, user)
	}
	rows.Close()

	for _, user := range users {
		result, err := repo.RawQuery(getUserTransactionSumQueryTemplate, user.Uuid, user.Uuid)
		if err != nil {
			log.Fatalln(err)
		}

		rows, ok := result.(*sql.Rows)
		if !ok {
			log.Fatalln("result is not rows")
		}

		var debSumNullable, valSumNullable sql.Null[int]
		var debSum, valSum int
		for rows.Next() {
			err := rows.Scan(&debSumNullable, &valSumNullable)
			if err != nil {
				log.Fatalln(err)
			}
			if debSumNullable.Valid {
				debSum = debSumNullable.V
			}
			if valSumNullable.Valid {
				valSum = valSumNullable.V
			}
		}
		rows.Close()

		finalBalance := user.Balance + valSum - debSum
		if finalBalance > 0 {
			trx := &model.Transaction{
				Uuid:     uuid.NewString(),
				Time:     time.Now(),
				Type:     model.PAYOFF,
				UserUuid: user.Uuid,
				TaskUuid: "-",
				Value:    finalBalance,
			}

			err := repo.Append(trx)
			if err != nil {
				log.Fatalln(err)
			}

			trxEventData := &eventmodel.TransactionEventData{
				Uuid:           trx.Uuid,
				Time:           trx.Time,
				Type:           string(trx.Type),
				LinkedUserUuid: trx.UserUuid,
				LinkedTaskUuid: "",
				Value:          trx.Value,
			}

			_, err = eventProducerTrxCUD.ProduceEventData(eventmodel.ACC_CUD_TRX_CREATED, trx.Uuid, reflect.TypeOf(*trx).String(), trxEventData, "1", nil)
			if err != nil {
				log.Fatalln(err)
			}

			// Генерим бизнес-событие по транзакции выплаты
			_, err = eventProducerTrxBE.ProduceEventData(eventmodel.ACC_BE_PAYOFF, trx.Uuid, reflect.TypeOf(*trx).String(), trxEventData, "1", nil)
			if err != nil {
				log.Fatalln(err)
			}

			user.Balance = 0

		} else {
			user.Balance = finalBalance
		}
		if finalBalance != 0 {
			err := repo.Update(user)
			if err != nil {
				log.Fatalln(err)
			}
		}

	}
}
