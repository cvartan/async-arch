// Взаимодействие с моделью данных сервиса авторизации

package main

import (
	model "async-arch/internal/domain/auth"
	"async-arch/internal/lib/base"
	repo "async-arch/internal/lib/base/repository/gorm"
	ou "async-arch/internal/lib/oshelper
	"log"
)

// Инициализация модели данных сервиса авторизации
func initModel() {

	host := ou.GetEnvValue("AUTHSRV_DBHOST", "192.168.1.99")
	dbname := "async-arch"
	user := "async-arch"
	password := "password"
	scheme := "auth"

	domainRepo, err := repo.CreateDomainRepository(host, dbname, scheme, user, password)
	if err != nil {
		log.Fatalln(err)
	}
	base.App.RegisterDomainRepository("auth", domainRepo)

	_, err = domainRepo.CreateObjectRepository(&model.User{})
	if err != nil {
		log.Fatal(err)
	}

}
