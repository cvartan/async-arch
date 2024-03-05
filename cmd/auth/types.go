// Определение специальных типов для методов-обработчиков запросов и событий

package main

import (
	events "async-arch/internal/lib/event"
)

// CreateUserRequest - данные запроса на добавление нового пользователя
type CreateUserRequest struct {
	Beak     string `json:"beak"`
	Name     string `json:"name"`
	Password string `json:"password"` // Пароль в запросе. Вместо пароля должен быть передан хэш MD5
	EMail    string `json:"email"`
	Role     string `json:"role"`
}

// UserResponse - структура используемая в ответах в которых возвращаются данные пользователя
type UserResponse struct {
	ID    uint   `json:"id"`
	Uuid  string `json:"uuid"`
	Beak  string `json:"beak"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// UserEventData - данные пользователя в событии
type UserEventData struct {
	Uuid  string `json:"uuid"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// Event - данные события для логирования в БД
type Event struct {
	ID uint `gorm:"PrimaryKey"`
	events.Event
	UserEventData
}

// CheckResponse - данные пользователя из проверенного JWT-токена
type CheckResponse struct {
	UserUuid string `json:"uuid"`
	UserRole string `json:"role"`
}
