package main

import (
	model "async-arch/internal/domain/accounting"
	"async-arch/internal/lib/base"
	events "async-arch/internal/lib/event"
	"encoding/json"
	"log"
	"strings"
)

func initEventHandlers() {
	eventConsumerCUD.AddConsumedEvent(events.AUTH_CUD_USER_CREATED, handleUserCreatedEvent)
	eventConsumerCUD.AddConsumedEvent(events.TASK_CUD_TASK_CREATED, handleTaskCreateEvent)
	eventConsumerCUD.Consume()
	eventConsumerBE.Consume()
}

// Обработчик CUD-события добавления нового пользователя
func handleUserCreatedEvent(event *events.Event, data interface{}) {
	// Получаем тело соообщения с событием
	body := string(data.([]byte))
	var userData UserEventData
	err := json.NewDecoder(strings.NewReader(body)).Decode(&userData)
	if err != nil {
		log.Fatal(err)
	}

	// Получаем данные пользователя и сохраняем их в БД
	user := &model.User{
		Uuid: userData.Uuid,
		Name: userData.Name,
	}

	repo, _ := base.App.GetDomainRepository("task")
	err = repo.Append(user)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Event %s(id=\"%s\") catched\n", event.EventType, event.EventID)
}

func handleTaskCreateEvent(event *events.Event, data interface{}) {
	// Получаем тело соообщения с событием
	body := string(data.([]byte))
	var taskData TaskEventData
	err := json.NewDecoder(strings.NewReader(body)).Decode(&taskData)
	if err != nil {
		log.Fatal(err)
	}

}
