package event

import "time"

// TaskEventData - данные задачи для CUD-события
type PricedTaskStreamData struct {
	TaskStreamData
	AssignmentTaskPrice int `json:"assignedPrice"`
	CompletedTaskPrice  int `json:"completedPrice"`
}

// TransactionEventData - данные транзакции для события
type TransactionEventData struct {
	Uuid           string    `json:"uuid"`
	Time           time.Time `json:"time"`
	Type           string    `json:"type"`
	LinkedUserUuid string    `json:"linkedUserUuid"`
	LinkedTaskUuid string    `json:"linkedTrxUuid"`
	Value          int       `json:"value"`
}
