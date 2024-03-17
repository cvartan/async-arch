package main

import (
	"async-arch/internal/lib/base"
	rabbitmq "async-arch/internal/lib/base/messages/rabbitmq"
	"async-arch/internal/lib/event"
	sysenv "async-arch/internal/lib/sysenv"
	"log"
)

var (
	eventProducerTaskCUD, eventProducerTrxCUD, eventProducerTrxBE *event.EventProducer
	eventConsumerCUD, eventConsumerBE                             *event.EventConsumer
)

func main() {
	// Запускаем http-сервер на порту 8091
	if err := base.App.InitHTTPServer("", 8092); err != nil {
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

	repo, _ := base.App.GetDomainRepository("accounting")

	eventProducerTrxBE, err = event.CreateEventProducer(manager, repo, "accounting", "account-lifecycle-log")
	if err != nil {
		log.Fatal(err)
	}

	eventProducerTaskCUD, err = event.CreateEventProducer(manager, repo, "accounting", "task-streaming")
	if err != nil {
		log.Fatal(err)
	}

	eventProducerTrxCUD, err = event.CreateEventProducer(manager, repo, "accounting", "trx-streaming")
	if err != nil {
		log.Fatal(err)
	}

	eventConsumerCUD = event.CreateEventConsumer(manager, repo, "acc-streaming")
	eventConsumerBE = event.CreateEventConsumer(manager, repo, "acc-business-events")

	// Так как ждать конца дня долго, то будем выполнять ребаланс по сообщению из очереди (куда сами сообщение и кидаем)
	err = base.App.Consume(manager.ID(), "rebalance-now", handleRebalanceMessage)
	if err != nil {
		log.Fatalln(err)
	}

	initEventHandlers()

	// Запускаем приложение
	base.App.Hold()
}
