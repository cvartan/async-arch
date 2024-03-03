package event

import (
	msg "async-arch/internal/lib/base/messages"
	str "async-arch/internal/lib/stringtool"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
)

// Используется для отправки событий в очередь
type EventProducer struct {
	producerId string
	manager    msg.MessageManager
	sender     string
}

// Создание продюсера событий
func CreateEventProducer(manager msg.MessageManager, sender, queueName string) (*EventProducer, error) {
	producerId := uuid.NewString()
	err := manager.AddProducer(producerId, queueName)
	if err != nil {
		return nil, err
	}
	return &EventProducer{
		producerId: producerId,
		manager:    manager,
		sender:     sender,
	}, nil
}

// Отправка события
func (p *EventProducer) ProduceEventData(eventType EventType, dataID, dataType string, data interface{}) (*Event, error) {
	evnt := Event{
		EventID:   uuid.NewString(),
		EventType: eventType,
		Subject:   dataType,
		Sender:    p.sender,
		CreatedAt: time.Now(),
	}
	headers := map[string]interface{}{
		"X_ID":        evnt.EventID,
		"X_EventType": string(evnt.EventType),
		"X_Subject":   evnt.Subject,
		"X_Sender":    evnt.Sender,
		"X_CreatedAt": evnt.CreatedAt.Format(time.RFC3339),
		"X-DataID":    dataID,
	}

	var dataJson string
	err := json.NewEncoder(str.CreateStringWriter(&dataJson)).Encode(data)
	if err != nil {
		return nil, err
	}

	producer, _ := p.manager.GetProducer(p.producerId)
	err = producer.ProduceMessage([]byte(dataID), []byte(dataJson), headers)
	if err != nil {
		return nil, err
	}

	return &evnt, nil
}

// Тип функции обработки события
type EventHandleFunc func(e *Event, data interface{})

// Обработчик событий
type EventConsumer struct {
	manager   msg.MessageManager
	queueName string
	events    map[EventType]EventHandleFunc
}

// Создание обработчика событий
func CreateEventConsumer(manager msg.MessageManager, queueName string) *EventConsumer {
	return &EventConsumer{
		manager:   manager,
		queueName: queueName,
		events:    make(map[EventType]EventHandleFunc),
	}
}

// Добавление обработчика для определенного типа события
func (p *EventConsumer) AddConsumedEvent(eventName EventType, handleFunc EventHandleFunc) {
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
				e.EventType = EventType(v.(string))
			}
		case "X_Subject":
			{
				e.Subject = v.(string)
			}
		case "X_Sender":
			{
				e.Sender = v.(string)
			}
		case "X_CreatedAt":
			{
				s := v.(string)
				e.CreatedAt, _ = time.Parse(time.RFC3339, s)
			}
		}
	}

	// Определяем обработчик для данного соьбытия и запускаем его
	if f, ok := p.events[e.EventType]; ok {
		f(e, value)
	}
}

// Запуск обработки событий
func (p *EventConsumer) Consume() {
	log.Printf("Start consuming events (manager = %s, queue = %s)\n", p.manager.ID(), p.queueName)
	p.manager.Consume(p.queueName, p.handleEvent)
}
