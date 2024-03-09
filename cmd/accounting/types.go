package main

import "time"

type GetLogResponse struct {
	Balance int `json:"balance"`
	Log     []TransactionLogItem
}

type TransactionLogItem struct {
	Uuid     string    `json:"uuid"`
	Time     time.Time `json:"time"`
	TaskUuid string    `json:"taskUuid"`
	Type     string    `json:"type"`
	Value    int       `json:"value"`
}
