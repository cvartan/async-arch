package base

import (
	msg "async-arch/lib/msg_srv"
	"errors"
)

// ServiceApplication - шаблон для сервиса
type ServiceApplication struct {
	messageManagers map[string]msg.MessageManager
}

// RegisterMessageManager - метод регистрации нового менеджера очередей
func (app *ServiceApplication) RegisterMessageManager(manager msg.MessageManager) {
	if _, ok := app.messageManagers[manager.ID()]; !ok {
		app.messageManagers[manager.ID()] = manager
	}
}

// RegisterMessageProducer - метод регистрации нового публикатора сообщений
func (app *ServiceApplication) RegisterMessageProducer(managerId, producerId string, queueName string) error {

	if managerId == "" {
		return errors.New("необходимо указать имя используемого менеджера очередей")
	}

	manager, ok := app.messageManagers[managerId]

	if !ok {
		return errors.New("менеджер с таким Id не найден")
	}

	if producerId == "" {
		return errors.New("необходимо указать имя-ключ для публикатора сообщений")
	}

	if queueName == "" {
		return errors.New("необходимо указать имя очереди для публикации")
	}

	if _, ok := manager.GetProducer(producerId); ok {
		return errors.New("публикатор с таким именем-ключом уже добавлен")
	}

	if err := manager.CreateProducer(producerId, queueName); err != nil {
		return err
	} else {
		return nil
	}
}

// SendMsg - метод отправки сообщений через укзанного публикатора сообщений
func (app *ServiceApplication) SendMsg(managerId, producerId string, key []byte, value []byte, headers map[string]interface{}) error {
	if managerId == "" {
		return errors.New("необходимо указать имя используемого менеджера очередей")
	}

	manager, ok := app.messageManagers[managerId]

	if !ok {
		return errors.New("менеджер с таким Id не найден")
	}

	if producerId == "" {
		return errors.New("необходимо указать имя-ключ для публикатора")
	}

	if len(key) == 0 {
		return errors.New("необходимо указать ключ для сообщения")
	}

	if producer, ok := manager.GetProducer(producerId); !ok {
		return errors.New("публикатор с таким именем не найден")
	} else {
		return producer.SendMsg(key, value, headers)
	}
}

// SendStrMsg - метод отправки сообщений в стоковом формате через указанного публикатора сообщений
func (app *ServiceApplication) SendStrMsg(managerId, producerId string, key string, value string, headers map[string]interface{}) error {
	return app.SendMsg(managerId, producerId, []byte(key), []byte(value), headers)
}

func (app *ServiceApplication) Close() {
	//1. Закрываем соединения для менеджеров очередей
	for _, manager := range app.messageManagers {
		manager.Close()
	}
	//2. Закрываем соединения для менеджеров баз данных
}

var App ServiceApplication

func init() {
	App.messageManagers = make(map[string]msg.MessageManager)
}
