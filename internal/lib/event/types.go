// Модель данных для работы с событиями
package event

import (
	"time"
)

type EventType string

const (
	TASK_BE_TASK_CREATED   EventType = "TASK.BE.TASK_CREATED"
	TASK_BE_TASK_ASSIGNED  EventType = "TASK.BE.TASK_ASSIGNED"
	TASK_BE_TASK_COMPLETED EventType = "TASK.BE.TASK_COMPLETED"

	AUTH_CUD_USER_CREATED      EventType = "AUTH.CUD.USER_CREATED"
	TASK_CUD_TASK_CREATED      EventType = "TASK.CUD.TASK_CREATED"
	TASK_CUD_TASK_UPDATED      EventType = "TASK.CUD.TASK_UPDATED"
	ACC_CUD_TASK_PRICED        EventType = "ACC.CUD.TASK_PRICED"
	ACC_CUD_TRX_CREATED        EventType = "ACC.CUD.TRX_CREATED"
	ACC_CUD_BALANCE_CALCULATED EventType = "ACC.CUD.BALANCE_CALCULATED"
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
