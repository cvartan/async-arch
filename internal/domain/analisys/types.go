package analisys

import "time"

// Структура для данных о пользователе
type User struct {
	ID   uint   `gorm:"primaryKey"`
	Uuid string `gorm:"unique"`
	Name string
}

// Структура для данных о задаче
type Task struct {
	ID             uint      `gorm:"primaryKey"`
	Uuid           string    `gorm:"unique"`
	IsComplete     bool      `gorm:"index"`
	CompleteTime   time.Time `gorm:"index"`
	CompletedPrice int
}

// Структура для данных о транзакции
type Transaction struct {
	ID       uint   `gorm:"primaryKey"`
	Uuid     string `gorm:"unique"`
	Type     string
	Time     time.Time
	UserUuid string `gorm:"not null;index"`
	Value    int
}
