package main

import (
	model "async-arch/internal/domain/analysis"
	"async-arch/internal/lib/base"
	repo "async-arch/internal/lib/base/repository/gorm"
	ou "async-arch/internal/lib/osutils"
	"log"
)

// Инициализация модели данных сервиса авторизации
func initModel() {

	host := ou.GetEnvValue("AUTHSRV_DBHOST", "192.168.1.99")
	dbname := "async-arch"
	user := "async-arch"
	password := "password"
	scheme := "analysis"

	domainRepo, err := repo.CreateDomainRepository(host, dbname, scheme, user, password)
	if err != nil {
		log.Fatalln(err)
	}
	base.App.RegisterDomainRepository("analysis", domainRepo)

	_, err = domainRepo.CreateObjectRepository(&model.Task{})
	if err != nil {
		log.Fatal(err)
	}

	_, err = domainRepo.CreateObjectRepository(&model.User{})
	if err != nil {
		log.Fatal(err)
	}

	_, err = domainRepo.CreateObjectRepository(&model.Transaction{})
	if err != nil {
		log.Fatal(err)
	}
}
