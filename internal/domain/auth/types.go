// Домен авторизации

package auth

import jwt "github.com/golang-jwt/jwt/v5"

// User - пользователь
type User struct {
	ID       uint       `gorm:"primaryKey"`
	Beak     ParrotBeak `gorm:"unique"`
	Password string     `gorm:"not null"`
	Uuid     string     `gorm:"not null"`
	Name     string     `gorm:"not null"`
	EMail    string     `gorm:"not null"`
	Role     UserRole   `gorm:"type:varchar(50);not null"`
}

// Уникальный профиль клюва попугая, используемый для аутентификации попугая
type ParrotBeak string

// Роли пользователей
type UserRole string

const (
	ADMIN     UserRole = "ADMIN"
	MANAGER   UserRole = "MANAGER"
	DEVELOPER UserRole = "DEVELOPER"
	ACCOUNTER UserRole = "ACCOUNTER"
)

type AuthClaims struct {
	jwt.RegisteredClaims
	UserUuid string
	UserRole string
}
