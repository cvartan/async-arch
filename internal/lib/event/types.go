// Модель данных для работы с событиями
package event

import (
	"time"
)

type EventType string

const (
	// Называем английским названием бизнес-события
	TASK_BE_TASK_CREATED   EventType = "TASK_CREATED"
	TASK_BE_TASK_ASSIGNED  EventType = "TASK_ASSIGNED"
	TASK_BE_TASK_COMPLETED EventType = "TASK_COMPLETED"
	// Называем по шаблону ДОМЕН.СУЩНОСТЬ.ВЫПОЛНЕННОЕ ДЕЙСТВИЕ
	AUTH_CUD_USER_CREATED      EventType = "AUTH.USER.CREATED"
	TASK_CUD_TASK_CREATED      EventType = "TASK.TASK.CREATED"
	TASK_CUD_TASK_UPDATED      EventType = "TASK.TASK.UPDATED"
	ACC_CUD_TASK_PRICED        EventType = "ACC.TASK.UPDATED"
	ACC_CUD_TRX_CREATED        EventType = "ACC.TRANSACTION.CREATED"
	ACC_CUD_BALANCE_CALCULATED EventType = "ACC.USER.UPDATED"
)

// Структура события
type Event struct {
	EventID   string    `json:"id"`
	EventType EventType `json:"type"`
	Subject   string    `json:"subject"`
	Sender    string    `json:"sender"`
	DataID    string    `json:"dataId"`
	CreatedAt time.Time `json:"createdAt"`
}

// Структура сообщения с событием
type EventMessage struct {
	Event
	Data interface{} `json:"data"`
}
