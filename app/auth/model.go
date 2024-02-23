package main

import (
	model "async-arch/model/domain/auth"
	"async-arch/util"
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
	connStr := util.GetEnvValue("AUTHSRV_DBCONNECT", "host=192.168.1.99 user=async-arch dbname=async-arch password=password sslmode=disable")
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
		model.TM_CREATE_TASK:           {"ALL"},
		model.TM_ASSIGN_TASK:           {model.ADMIN, model.MANAGER},
		model.TM_VIEW_SELF_TASKS:       {"ALL"},
		model.TM_VIEW_ALL_TASKS:        {model.ADMIN},
		model.ACC_VIEW_BALANCE:         {model.ADMIN, model.ACCOUNTER},
		model.ACC_VIEW_SELF_BALANCE:    {"ALL"},
		model.STAT_VIEW_PRICING_INFO:   {model.ADMIN},
		model.STAT_VIEW_BALANCE_BY_DAY: {model.ADMIN},
	}
}
