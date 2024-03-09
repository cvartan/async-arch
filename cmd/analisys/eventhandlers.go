package main

import (
	model "async-arch/internal/domain/analisys"
	eventmodel "async-arch/internal/domain/event"
	"async-arch/internal/lib/base"
	events "async-arch/internal/lib/event"
	"encoding/json"
	"log"
	"strings"
)

func initEventHandlers() {
	eventConsumerCUD.AddConsumedEvent(eventmodel.AUTH_CUD_USER_CREATED, "1", handleUserCreateEvent)
	eventConsumerCUD.AddConsumedEvent(eventmodel.TASK_CUD_TASK_CREATED, "1", handleTaskCreateEvent)
	eventConsumerCUD.AddConsumedEvent(eventmodel.TASK_CUD_TASK_UPDATED, "1", handleTaskUpdateEvent)
	eventConsumerCUD.AddConsumedEvent(eventmodel.ACC_CUD_TASK_PRICED, "1", handleTaskUpdateEvent)
	eventConsumerCUD.AddConsumedEvent(eventmodel.ACC_CUD_TRX_CREATED, "1", handleTransactionCreateEvent)
	eventConsumerCUD.Consume()
}

// Обработка события добавления нового пользователя
func handleUserCreateEvent(event *events.Event, data interface{}) {
	log.Printf("Event %s(id=\"%s\",dataid=\"%s\") received\n", event.EventType, event.EventID, event.DataID)

	userData := &eventmodel.UserEventData{}
	err := json.NewDecoder(strings.NewReader(string(data.([]byte)))).Decode(userData)
	if err != nil {
		log.Fatalln(err)
	}

	user := &model.User{
		Uuid: userData.Uuid,
		Name: userData.Name,
	}

	repo, _ := base.App.GetDomainRepository("analisys")
	err = repo.Append(user)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Event %s(id=\"%s\") catched\n", event.EventType, event.EventID)
}

// Обработка события создания новой задачи
func handleTaskCreateEvent(event *events.Event, data interface{}) {
	log.Printf("Event %s(id=\"%s\",dataid=\"%s\") received\n", event.EventType, event.EventID, event.DataID)

	taskData := &eventmodel.TaskEventData{}
	err := json.NewDecoder(strings.NewReader(string(data.([]byte)))).Decode(taskData)
	if err != nil {
		log.Fatalln(err)
	}

	task := &model.Task{
		Uuid: taskData.Uuid,
	}

	repo, _ := base.App.GetDomainRepository("analisys")
	err = repo.Append(task)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Event %s(id=\"%s\") catched\n", event.EventType, event.EventID)
}

// Обработка события обновления задачи
// Также тут происходит обновление события формируемого в сервисе аккаунтинга после назначения цен
func handleTaskUpdateEvent(event *events.Event, data interface{}) {
	log.Printf("Event %s(id=\"%s\",dataid=\"%s\") received\n", event.EventType, event.EventID, event.DataID)

	taskData := &eventmodel.PricedTaskStreamData{}
	err := json.NewDecoder(strings.NewReader(string(data.([]byte)))).Decode(taskData)
	if err != nil {
		log.Fatalln(err)
	}

	task := &model.Task{}

	repo, _ := base.App.GetDomainRepository("analisys")
	err = repo.Get(task, map[string]interface{}{"uuid": taskData.Uuid})
	if err != nil {
		log.Fatalln(err)
	}

	if taskData.CompletedTaskPrice != 0 {
		task.CompletedPrice = taskData.CompletedTaskPrice
	}

	if taskData.State == "COMPLETED" {
		task.IsComplete = true
		task.CompleteTime = event.CreatedAt
	}

	err = repo.Update(task)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Event %s(id=\"%s\") catched\n", event.EventType, event.EventID)
}

// Обработка события создания новой транзакции
func handleTransactionCreateEvent(event *events.Event, data interface{}) {
	log.Printf("Event %s(id=\"%s\",dataid=\"%s\") received\n", event.EventType, event.EventID, event.DataID)

	trxData := &eventmodel.TransactionEventData{}
	err := json.NewDecoder(strings.NewReader(string(data.([]byte)))).Decode(trxData)
	if err != nil {
		log.Fatalln(err)
	}

	trx := &model.Transaction{
		Uuid:     trxData.Uuid,
		Type:     trxData.Type,
		Time:     trxData.Time,
		UserUuid: trxData.LinkedUserUuid,
		Value:    trxData.Value,
	}

	repo, _ := base.App.GetDomainRepository("analisys")
	err = repo.Append(trx)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Event %s(id=\"%s\") catched\n", event.EventType, event.EventID)
}
