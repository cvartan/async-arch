// Модель данных для работы с событиями
package event

import (
	model "async-arch/internal/domain/event"
	"time"
)

// Структура события
type Event struct {
	EventID   string          `json:"id"`
	EventType model.EventType `json:"type"`
	Subject   string          `json:"subject"`
	Sender    string          `json:"sender"`
	DataID    string          `json:"dataId"`
	CreatedAt time.Time       `json:"createdAt"`
	Version   string          `json:"version"`
}

// Структура сообщения с событием
type EventMessage struct {
	Event
	Data interface{} `json:"data"`
}

// Структура записи с событием в БД
type EventLog struct {
	ID uint `gorm:"primaryKey"`
	Event
	Source string
}

type MessageHeader struct {
	Header string `json:"header"`
	Value  string `json:"value"`
}

type DeadEventSource struct {
	Headers []*MessageHeader `json:"headers"`
	Source  string           `json:"source"`
}

type DeadEvent struct {
	ID          uint `gorm:"primaryKey"`
	MessageKey  string
	MessageBody string
	ErrorText   string
}
