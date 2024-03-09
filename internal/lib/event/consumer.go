package event

import (
	model "async-arch/internal/domain/event"
	msg "async-arch/internal/lib/base/messages"
	repo "async-arch/internal/lib/base/repository"
	"async-arch/internal/lib/schema"
	"async-arch/internal/lib/stringtool"
	"errors"

	"encoding/base64"
	"encoding/json"
	"log"
	"time"
)

// Тип функции обработки события
type EventHandleFunc func(e *Event, data interface{})

// Обработчик событий
type EventConsumer struct {
	manager       msg.MessageManager
	repository    repo.DomainRepositoryManager
	validator     *schema.SchemaValidator
	queueName     string
	events        map[model.EventType]EventHandleFunc
	eventVersions map[model.EventType]string
}

// Создание обработчика событий
func CreateEventConsumer(
	manager msg.MessageManager,
	repository repo.DomainRepositoryManager,
	queueName string,
) *EventConsumer {
	if repository != nil {
		_, err := repository.CreateObjectRepository(&DeadEvent{})
		if err != nil {
			return nil
		}
	}

	return &EventConsumer{
		manager:       manager,
		repository:    repository,
		validator:     schema.CreateSchemaValidator(),
		queueName:     queueName,
		events:        make(map[model.EventType]EventHandleFunc),
		eventVersions: make(map[model.EventType]string),
	}
}

// Добавление обработчика для определенного типа события
func (p *EventConsumer) AddConsumedEvent(eventName model.EventType, eventVersion string, handleFunc EventHandleFunc) {
	p.events[eventName] = handleFunc
	p.eventVersions[eventName] = eventVersion
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
		p.saveDeadLetter(key, value, headers, errors.New("incompatible event"))
		return
	}

	// Выполняем валидацию события
	version := p.eventVersions[e.EventType]

	body := string(value)
	err := p.validator.Validate(string(e.EventType), version, body)
	if err != nil && p.repository != nil {
		// Сохраняем сообщение в таблицу "мертвых сообщений"
		p.saveDeadLetter(key, value, headers, err)
		return
	}

	// Запускаем обработку события
	f(e, value)

}

func (p *EventConsumer) saveDeadLetter(key, value []byte, headers map[string]interface{}, err error) {
	deadEventSource := &DeadEventSource{}
	var messageHeaders []*MessageHeader
	for k, v := range headers {
		h := &MessageHeader{
			Header: k,
			Value:  v.(string),
		}
		messageHeaders = append(messageHeaders, h)
	}
	deadEventSource.Headers = messageHeaders
	deadEventSource.Source = base64.StdEncoding.EncodeToString(value)

	var deadEventJson string
	err1 := json.NewEncoder(stringtool.CreateStringWriter(&deadEventJson)).Encode(deadEventSource)
	if err1 != nil {
		log.Println(err1)
		return
	}

	deadEvent := &DeadEvent{
		MessageKey:  string(key),
		MessageBody: deadEventJson,
		ErrorText:   err.Error(),
	}

	err1 = p.repository.Append(deadEvent)
	if err != nil {
		log.Println(err1)
	}
}

// Запуск обработки событий
func (p *EventConsumer) Consume() {
	log.Printf("Start consuming events (manager = %s, queue = %s)\n", p.manager.ID(), p.queueName)
	p.manager.Consume(p.queueName, p.handleEvent)
}
