// Домен управления задачами

package taskman

import (
	authmodel "async-arch/internal/domain/auth"
)

// Данные пользователя
type User struct {
	ID   uint   `gorm:"primaryKey"`
	Uuid string `gorm:"unique"`
	Name string `gorm:"not null"`
	Role authmodel.UserRole
}

// Задача
type Task struct {
	ID   uint   `gorm:"primaryKey"`
	Uuid string `gorm:"unique"`
	// Расширили модель двумя новыми атрибутами
	Title            string
	JiraId           string
	Description      string `gorm:"not null"`
	AssignedUserUuid string `gorm:"not null"`
	State            TaskState
}

// Состояния задачи
type TaskState string

const (
	TASK_ACTIVE    TaskState = "ACTIVE"
	TASK_COMPLETED TaskState = "COMPLETED"
)
