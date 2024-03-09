// Типы данных используемые в обработчиках запросов или событий

package main

import (
	model "async-arch/internal/domain/taskman"
)

// CreateTaskRequest - данные запроса создания задачи
type CreateTaskRequest struct {
	Description string `json:"description"`
}

// TaskResponse - данные задачи, возвращаемые в ответах
type TaskResponse struct {
	ID          uint            `json:"id"`
	Uuid        string          `json:"uuid"`
	Description string          `json:"description"`
	UserUuid    string          `json:"userUuid"`
	State       model.TaskState `json:"state"`
}

// ReassingTasksResponse - список переназначенных задач
type ReassingTasksResponse []ReassignTasksResponseItem

type ReassignTasksResponseItem struct {
	ID       uint   `json:"id"`
	Uuid     string `json:"uuid"`
	UserUuid string `json:"userUuid"`
}

// UserTasksList - список задач
type UserTasksList []TaskResponse
