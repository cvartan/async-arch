package main

import (
	authlib "async-arch/internal/lib/auth"
	model "async-arch/internal/model/domain/auth"
	"async-arch/internal/sysenv"
	"log"
	"time"

	psql "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var db *gorm.DB

type sessionParams struct {
	UserID  uint
	Expires time.Time
}

var sessions map[string]sessionParams

type roles []model.UserRole

var permissions map[string]roles

func initModel() {
	connStr := sysenv.GetEnvValue("AUTHSRV_DBCONNECT", "host=192.168.1.99 user=async-arch dbname=async-arch password=password sslmode=disable")
	var err error
	db, err = gorm.Open(psql.Open(connStr), &gorm.Config{NamingStrategy: schema.NamingStrategy{TablePrefix: "auth.", SingularTable: true}})
	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&model.User{})
}

func init() {
	sessions = make(map[string]sessionParams)

	permissions = map[string]roles{
		authlib.TM_CREATE_TASK:           {"ALL"},
		authlib.TM_ASSIGN_TASK:           {model.ADMIN, model.MANAGER},
		authlib.TM_VIEW_SELF_TASKS:       {"ALL"},
		authlib.TM_VIEW_ALL_TASKS:        {model.ADMIN},
		authlib.ACC_VIEW_BALANCE:         {model.ADMIN, model.ACCOUNTER},
		authlib.ACC_VIEW_SELF_BALANCE:    {"ALL"},
		authlib.STAT_VIEW_PRICING_INFO:   {model.ADMIN},
		authlib.STAT_VIEW_BALANCE_BY_DAY: {model.ADMIN},
	}
}
