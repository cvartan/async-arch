package main

import (
	"async-arch/internal/lib/base"
	rabbitmq "async-arch/internal/lib/base/messages/rabbitmq"
	"async-arch/internal/lib/event"
	sysenv "async-arch/internal/lib/sysenv"
	"log"
)

var (
	eventConsumerCUD, eventConsumerBE *event.EventConsumer
)

func main() {
	// Запускаем http-сервер на порту 8093
	if err := base.App.InitHTTPServer("", 8093); err != nil {
		log.Fatalln(err)
	}

	// Инициализируем модель данных в БД
	initModel()
	// Инициализируем обработчики запросов http
	initHandlers()
	// Инициализируем менеджер очередей и продюсер для событий
	server := sysenv.GetEnvValue("RABBITMQ_SERVER", "192.168.1.99")
	vhostName := sysenv.GetEnvValue("RABBITMQ_VHOST", "async_arch")
	user := "asyncarch"
	password := "password"
	manager, err := rabbitmq.CreateRabbitMQManager(
		"rabbit_mq",
		user,
		password,
		server,
		vhostName,
		rabbitmq.DEFAULT_PORT,
	)
	if err != nil {
		log.Fatal(err)
	}
	base.App.RegisterMessageManager(manager)

	repo, _ := base.App.GetDomainRepository("analisys")
	eventConsumerCUD = event.CreateEventConsumer(manager, repo, "analisys-streaming")
	eventConsumerBE = event.CreateEventConsumer(manager, repo, "analysis-bussiness-events")

	initEventHandlers()

	// Запускаем приложение
	base.App.Hold()
}
