package base

import (
	msg "async-arch/lib/msg_srv"
	"errors"
)

// ServiceApplication - шаблон для сервиса
type ServiceApplication struct {
	messageProducers map[string]msg.MessageProducer
	messageConsumers map[string]msg.MessageConsumer
}

// RegisterMessageProducer - метод регистрации нового публикатора сообщений
func (app *ServiceApplication) RegisterMessageProducer(name string, producer msg.MessageProducer) error {
	if name == "" {
		return errors.New("необходимо указать имя-ключ для публикатора сообщений")
	}

	if producer == nil {
		return errors.New("необходимо указать используемый публикатор сообщений")
	}
	_, ok := app.messageProducers[name]

	if ok {
		return errors.New("публикатор с таким именем-ключом уже добавлен")
	}

	app.messageProducers[name] = producer

	return nil
}

// RegisterMessageConsumer - метод регистрации нового читатетля сообщений
func (app *ServiceApplication) RegisterMessageConsumer(name string, consumer msg.MessageConsumer) error {
	if name == "" {
		return errors.New("необходимо указать имя-ключ для читателя сообщений")
	}

	if consumer == nil {
		return errors.New("необходимо указать используемый читатель сообщений")
	}
	_, ok := app.messageConsumers[name]

	if ok {
		return errors.New("читатель с таким именем-ключом уже добавлен")
	}

	app.messageConsumers[name] = consumer

	return nil
}

// SendMsg - метод отправки сообщений через укзанного публикатора сообщений
func (app *ServiceApplication) SendMsg(name string, key []byte, value []byte, headers []string) error {
	if name == "" {
		return errors.New("необходимо указать имя-ключ для публикатора")
	}

	if len(key) == 0 {
		return errors.New("необходимо указать ключ для сообщения")
	}

	if producer, ok := app.messageProducers[name]; !ok {
		return errors.New("публикатор с таким именем не найден")
	} else {
		return producer.SendMsg(key, value, headers)
	}
}

// SendStrMsg - метод отправки сообщений в стоковом формате через указанного публикатора сообщений
func (app *ServiceApplication) SendStrMsg(name string, key string, value string, headers []string) error {
	return app.SendMsg(name, []byte(key), []byte(value), headers)
}

// ReadMsg -метод запуска чтения сообщений через указанного читателя
func (app *ServiceApplication) ReadMsg(name string, msgHandler msg.MessageHandler) error {
	if name == "" {
		return errors.New("необходимо указать имя-ключ для читателя")
	}

	if msgHandler == nil {
		return errors.New("необходимо указать функцию обработки сообщения")
	}

	if consumer, ok := app.messageConsumers[name]; !ok {
		return errors.New("читатель с таким именем не найден")
	} else {
		return consumer.ReadMsg(msgHandler)
	}

}

func (app *ServiceApplication) Close() {
	//1. Закрываем соединения для менеджеров очередей
	for _, consumer := range app.messageConsumers {
		consumer.Close()
	}
	//2. Закрываем соединения для менеджеров баз данных
}

var App ServiceApplication

func init() {
	App.messageProducers = make(map[string]msg.MessageProducer)
	App.messageConsumers = make(map[string]msg.MessageConsumer)
}
