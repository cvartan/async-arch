package event

import (
	model "async-arch/internal/domain/event"
	msg "async-arch/internal/lib/base/messages"
	repo "async-arch/internal/lib/base/repository"

	"log"
	"time"
)

// Тип функции обработки события
type EventHandleFunc func(e *Event, data interface{})

// Обработчик событий
type EventConsumer struct {
	manager    msg.MessageManager
	repository repo.DomainRepositoryManager
	queueName  string
	events     map[model.EventType]EventHandleFunc
}

// Создание обработчика событий
func CreateEventConsumer(
	manager msg.MessageManager,
	repository repo.DomainRepositoryManager,
	queueName string,
) *EventConsumer {
	return &EventConsumer{
		manager:    manager,
		repository: repository,
		queueName:  queueName,
		events:     make(map[model.EventType]EventHandleFunc),
	}
}

// Добавление обработчика для определенного типа события
func (p *EventConsumer) AddConsumedEvent(eventName model.EventType, eventVersion string, handleFunc EventHandleFunc) {
	p.events[eventName] = handleFunc
}

// Функция обработки сообщения с событием
func (p *EventConsumer) handleEvent(key, value []byte, headers map[string]interface{}) {
	e := &Event{}
	for k, v := range headers {
		switch k {
		case "X_ID":
			{
				e.EventID = v.(string)
			}
		case "X_EventType":
			{
				e.EventType = model.EventType(v.(string))
			}
		case "X_Subject":
			{
				e.Subject = v.(string)
			}
		case "X_Sender":
			{
				e.Sender = v.(string)
			}
		case "X_DataID":
			{
				e.DataID = v.(string)
			}
		case "X_CreatedAt":
			{
				s := v.(string)
				e.CreatedAt, _ = time.Parse(time.RFC3339, s)
			}
		case "X_Version":
			{
				e.Version = v.(string)
			}
		}
	}
	// Определяем обработчик для данного события
	f, ok := p.events[e.EventType]
	if !ok {
		return
	}

	// Запускаем обработку события
	f(e, value)

}

// Запуск обработки событий
func (p *EventConsumer) Consume() {
	log.Printf("Start consuming events (manager = %s, queue = %s)\n", p.manager.ID(), p.queueName)
	p.manager.Consume(p.queueName, p.handleEvent)
}
