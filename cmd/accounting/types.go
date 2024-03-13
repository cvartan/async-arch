package main

import "time"

// Ответ по статистике для пользовтаеля
type GetLogResponse struct {
	Balance int                  `json:"balance"` // текущий баланс
	Log     []TransactionLogItem // лог операций
}

// Запись в логе операций
type TransactionLogItem struct {
	Uuid     string    `json:"uuid"`
	Time     time.Time `json:"time"`
	TaskUuid string    `json:"taskUuid"`
	Type     string    `json:"type"`
	Value    int       `json:"value"`
}
