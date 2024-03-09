// Обработчики событий

package main

import (
	auth "async-arch/internal/domain/auth"
	eventmodel "async-arch/internal/domain/event"
	model "async-arch/internal/domain/taskman"
	base "async-arch/internal/lib/base"
	events "async-arch/internal/lib/event"
	"encoding/json"
	"log"
	"strings"
)

// UserEventData - данные пользователя в событии
type UserEventData struct {
	Uuid  string `json:"uuid"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// Инициализация обработчиков по событиям
func initEventHandlers() {
	eventConsumer.AddConsumedEvent(eventmodel.AUTH_CUD_USER_CREATED, "1", handleUserCreatedEvent)
	eventConsumer.Consume()
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
		Role: auth.UserRole(userData.Role),
	}

	repo, _ := base.App.GetDomainRepository("task")
	err = repo.Append(user)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Event %s(id=\"%s\") catched\n", event.EventType, event.EventID)
}
