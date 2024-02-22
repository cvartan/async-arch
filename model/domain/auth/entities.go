package auth

import "gorm.io/gorm"

// User - пользователь
type User struct {
	gorm.Model
	ParrotBeak  string           `gorm:"unique"`
	Name        string           `gorm:"Name;not null"`
	Role        UserRole         `gorm:"type:varchar(50);not null"`
	Permissions []UserPermission `gorm:"foreignKey:UserID"`
}

type UserRole string

const (
	ADMIN     UserRole = "ADMIN"
	MANAGER   UserRole = "MANAGER"
	DEVELOPER UserRole = "DEVELOPER"
	ACCOUNTER UserRole = "ACCOUNTER"
)

func CreateNewUser(ParrotBeak, Name, Role string) *User {
	return &User{
		ParrotBeak: ParrotBeak,
		Name:       Name,
		Role:       UserRole(Role),
	}
}

type UserPermission struct {
	gorm.Model
	UserID     uint
	Permission string `gorm:"type:varchar(50)"`
}
