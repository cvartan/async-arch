package main

import (
	model "async-arch/model/domain/auth"
	"async-arch/util"
	"log"

	psql "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var db *gorm.DB

func initModel() {
	connStr := util.GetEnvValue("AUTHSRV_DBCONNECT", "host=192.168.1.99 user=async-arch dbname=async-arch password=password sslmode=disable")
	var err error
	db, err = gorm.Open(psql.Open(connStr), &gorm.Config{NamingStrategy: schema.NamingStrategy{TablePrefix: "auth.", SingularTable: true}})
	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&model.User{}, &model.UserPermission{})
}
