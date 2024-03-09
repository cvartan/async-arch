package event

import (
	model "async-arch/internal/domain/event"
	msg "async-arch/internal/lib/base/messages"
	repo "async-arch/internal/lib/base/repository"
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
	repository repo.DomainRepositoryManager
	sender     string
}

// Создание продюсера событий
func CreateEventProducer(
	manager msg.MessageManager, // Менеджер сообщений
	repository repo.DomainRepositoryManager, // Доменный репозиторий для сохранения лога событий в БД
	sender, // Имя сервиса - отправителя сообщений
	queueName string, // Имя очереди в которую помещаются исходящие сообщения
) (
	*EventProducer, // Продюсер событий
	error, // Ошибка
) {
	producerId := uuid.NewString()
	err := manager.AddProducer(producerId, queueName)
	if err != nil {
		return nil, err
	}
	if repository != nil {
		// Создаем репозиторий для события и заодно создаем таблицу в БД, если она отсутствует
		_, err := repository.CreateObjectRepository(&EventLog{})
		if err != nil {
			return nil, err
		}
	}

	return &EventProducer{
		producerId: producerId,
		manager:    manager,
		repository: repository,
		sender:     sender,
	}, nil
}

// Отправка события
func (p *EventProducer) ProduceEventData(eventType model.EventType, dataID, dataType string, data interface{}, eventVersion string) (*Event, error) {
	evnt := Event{
		EventID:   uuid.NewString(),
		EventType: eventType,
		Subject:   dataType,
		Sender:    p.sender,
		CreatedAt: time.Now(),
		Version:   eventVersion,
	}
	headers := map[string]interface{}{
		"X_ID":        evnt.EventID,
		"X_EventType": string(evnt.EventType),
		"X_Subject":   evnt.Subject,
		"X_Sender":    evnt.Sender,
		"X_CreatedAt": evnt.CreatedAt.Format(time.RFC3339),
		"X_DataID":    dataID,
		"X_Version":   evnt.Version,
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

	// Сохраняем событие в БД (если репозиторий определен)
	if p.repository != nil {
		func() {
			var eventSourceText string

			err := json.NewEncoder(str.CreateStringWriter(&eventSourceText)).Encode(data)

			if err != nil {
				log.Println(err)
				return
			}

			el := &EventLog{
				Event:  evnt,
				Source: eventSourceText,
			}

			err = p.repository.Append(el)
			if err != nil {
				log.Println(err)
			}
		}()
	}

	return &evnt, nil
}
