package main

import (
	"async-arch/internal/lib/base"
	rabbitmq "async-arch/internal/lib/base/messages/rabbitmq"
	"async-arch/internal/lib/event"
	sysenv "async-arch/internal/lib/sysenv"
	"log"
)

var (
	eventProducerCUD, eventProducerBE *event.EventProducer
	eventConsumer                     *event.EventConsumer
)

func main() {
	// Запускаем http-сервер на порту 8091
	if err := base.App.InitHTTPServer("", 8091); err != nil {
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
	manager, err := rabbitmq.CreateRabbitMQManager("rabbit_mq", user, password, server, vhostName, rabbitmq.DEFAULT_PORT)
	if err != nil {
		log.Fatal(err)
	}
	base.App.RegisterMessageManager(manager)

	eventProducerCUD, err = event.CreateEventProducer(manager, "task-streaming", "task-streaming")
	if err != nil {
		log.Fatal(err)
	}

	eventProducerBE, err = event.CreateEventProducer(manager, "task-lifecycle", "task-lifecycle-log")
	if err != nil {
		log.Fatal(err)
	}

	eventConsumer = event.CreateEventConsumer(manager, "task-streaming")
	initEventHandlers()

	// Запускаем приложение
	base.App.Hold()
}
