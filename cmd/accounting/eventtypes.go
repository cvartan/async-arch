package main

import (
	taskmodel "async-arch/internal/domain/taskman"
)

// UserEventData - данные пользователя в событии
type UserEventData struct {
	Uuid  string `json:"uuid"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// TaskEventData - данные задачи для события
type TaskEventData struct {
	Uuid             string              `json:"uuid"`
	Description      string              `json:"description"`
	AssignedUserUuid string              `json:"assignedUserUuid"`
	State            taskmodel.TaskState `json:"state"`
}

// TransactionEventData - данные транзакции для события
type TransactionEventData struct {
	Uuid           string `json:"uuid"`
	LinkedUserUuid string `json:"linkedUserUuid"`
	LinkedTaskUuid string `json:"linkedTrxUuid"`
}
