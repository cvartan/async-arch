package auth

import "gorm.io/gorm"

// User - пользователь
type User struct {
	gorm.Model
	ParrotBeak string   `gorm:"unique"`
	Name       string   `gorm:"not null"`
	EMail      string   `gorm:"not null"`
	Role       UserRole `gorm:"type:varchar(50);not null"`
}

type UserRole string

const (
	ADMIN     UserRole = "ADMIN"
	MANAGER   UserRole = "MANAGER"
	DEVELOPER UserRole = "DEVELOPER"
	ACCOUNTER UserRole = "ACCOUNTER"
)

func CreateNewUser(ParrotBeak, Name, EMail, Role string) *User {
	return &User{
		ParrotBeak: ParrotBeak,
		Name:       Name,
		EMail:      EMail,
		Role:       UserRole(Role),
	}
}
