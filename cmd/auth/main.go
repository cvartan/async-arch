package main

import (
	"async-arch/internal/lib/base"
	rabbitmq "async-arch/internal/lib/base/messages/rabbitmq"
	"async-arch/internal/lib/event"
	ou "async-arch/internal/lib/osutils"
	"log"
)

var eventProducer *event.EventProducer

func main() {
	// Запускаем http-сервер на порту 8090
	if err := base.App.InitHTTPServer("", 8090); err != nil {
		log.Fatalln(err)
	}

	// Инициализируем модель данных в БД
	initModel()
	// Инициализируем обработчики запросов http
	initHandlers()
	// Инициализируем менеджер очередей и продюсер для событий
	server := ou.GetEnvValue("RABBITMQ_SERVER", "192.168.1.99")
	vhostName := ou.GetEnvValue("RABBITMQ_VHOST", "async_arch")
	user := "asyncarch"
	password := "password"
	manager, err := rabbitmq.CreateRabbitMQManager("rabbit_mq", user, password, server, vhostName, rabbitmq.DEFAULT_PORT)
	if err != nil {
		log.Fatal(err)
	}
	base.App.RegisterMessageManager(manager)

	// Инициализируем продюсер событий
	repo, _ := base.App.GetDomainRepository("auth")
	eventProducer, err = event.CreateEventProducer(manager, repo, "auth", "user-streaming")
	if err != nil {
		log.Fatal(err)
	}

	// Фиксируем запуск приложения
	base.App.Hold()
}
