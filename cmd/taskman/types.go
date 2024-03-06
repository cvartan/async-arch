// Типы данных используемые в обработчиках запросов или событий

package main

import (
	model "async-arch/internal/domain/taskman"
	events "async-arch/internal/lib/event"
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

// Event - данные события для сохранения в лог событий
type Event struct {
	ID uint `gorm:"PrimaryKey"`
	events.Event
	TaskEventData
}

// TaskEventData - данные задачи для события
type TaskEventData struct {
	Uuid             string `json:"uuid"`
	Description      string `json:"description"`
	AssignedUserUuid string `json:"userUuid"`
}

// UserTasksList - список задач
type UserTasksList []TaskResponse
