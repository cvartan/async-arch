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
	ID             uint   `gorm:"primaryKey"`
	Uuid           string `gorm:"unique"`
	Title          string
	JiraId         string
	Description    string
	IsComplete     bool      `gorm:"index"`
	CompleteTime   time.Time `gorm:"index"`
	AssignedPrice  int
	CompletedPrice int
}

// Структура для данных о транзакции
type Transaction struct {
	ID       uint   `gorm:"primaryKey"`
	Uuid     string `gorm:"unique"`
	Type     string
	Time     time.Time
	UserUuid string `gorm:"not null;index"`
	TaskUuid string
	Value    int
}

// Структура для хранения бизнес-транзакций
type BusinessEvent struct {
	ID        uint   `gorm:"primaryKey"`
	Uuid      string `gorm:"unique"`
	EventType string
	DataType  string
	DataUuid  string
	Time      time.Time
	Sender    string
	Data      string
}
