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

const (
	TM_CREATE_TASK           string = "TM.CREATE_TASK"
	TM_ASSIGN_TASK           string = "TM.ASSIGN_TASK"
	TM_VIEW_SELF_TASKS       string = "TM.VIEW_SELF_TASKS"
	TM_VIEW_ALL_TASKS        string = "TN.VIEW_ALL_TASKS"
	ACC_VIEW_BALANCE         string = "ACC.VIEW_BALANCE"
	ACC_VIEW_SELF_BALANCE    string = "ACC.VIEW_SELF_BALANCE"
	STAT_VIEW_PRICING_INFO   string = "STAT.VIEW_PRICING_INFO"
	STAT_VIEW_BALANCE_BY_DAY string = "STAT.VIEW_BALANCE_BY_DAY"
)
