// Домен аккаунтинга

package accounting

import "time"

// Данные пользователя
type User struct {
	ID      uint    `gorm:"primaryKey"`
	Uuid    string  `gorm:"unique"`
	Name    string  `gorm:"not null"`
	Balance float32 `gorm:"default:0"`
}

// Данные задачи
type Task struct {
	ID                  uint    `gorm:"primaryKey"`
	Uuid                string  `gorm:"unique"`
	AssignmentTaskPrice float32 `gorm:"not null"`
	CompleteTaskPrice   float32 `gorm:"not null"`
}

// Транзакция
type Transaction struct {
	ID       uint            `gorm:"primaryKey"`
	Time     time.Time       `gorm:"not null"`
	UserUuid string          `gorm:"not null"`
	TaskUuid string          `gorm:"not null"`
	Type     TransactionType `gorm:"not null"`
	Value    float32         `gorm:"not null"`
}

// Тип транзакции
type TransactionType string

const (
	DEBITING TransactionType = "DEBITING" // Списание
	VALUE    TransactionType = "VALUE"    // Зачисление
	PAYOFF   TransactionType = "PAYOFF"   // Выплата
)
