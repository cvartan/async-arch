package base

import (
	msg "async-arch/lib/msgconnect"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/xlab/closer"
)

// serviceApplication - шаблон для сервиса
type serviceApplication struct {
	messageManagers map[string]msg.MessageManager
	httpServer      *http.Server
}

// RegisterMessageManager - метод регистрации нового менеджера очередей
func (app *serviceApplication) RegisterMessageManager(manager msg.MessageManager) {
	if _, ok := app.messageManagers[manager.ID()]; !ok {
		app.messageManagers[manager.ID()] = manager
	}
}

// RegisterMessageProducer - метод регистрации нового публикатора сообщений
func (app *serviceApplication) RegisterMessageProducer(managerId, producerId string, queueName string) error {

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
func (app *serviceApplication) SendMsg(managerId, producerId string, key []byte, value []byte, headers map[string]interface{}) error {
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

// SendStrMsg - метод отправки сообщений в строковом формате через указанного публикатора сообщений
func (app *serviceApplication) SendStrMsg(managerId, producerId string, key string, value string, headers map[string]interface{}) error {
	return app.SendMsg(managerId, producerId, []byte(key), []byte(value), headers)
}

// Consume - метод прослушивания очереди
func (app *serviceApplication) Consume(managerId, queueName string, handler msg.MessageHandler) error {
	if managerId == "" {
		return errors.New("необходимо указать имя используемого менеджера очередей")
	}

	manager, ok := app.messageManagers[managerId]
	if !ok {
		return errors.New("менеджер с таким Id не найден")
	}

	return manager.Consume(queueName, handler)

}

// InitHTTPServer - инициализация http сервера (запуск сервера выполняется при старте приложения)
func (app *serviceApplication) InitHTTPServer(address string, port uint16) error {
	if app.httpServer != nil {
		return errors.New("http-сервер уже инициализирован")
	}
	if port == 0 {
		port = 80
	}
	completeAddress := fmt.Sprintf("%s:%d", address, port)

	mux := http.NewServeMux()

	app.httpServer = &http.Server{
		Addr:    completeAddress,
		Handler: mux,
	}

	return nil
}

// HandleFunc - назначение метода обработки шаблону пути (пути с параметрами требуют версии go от 1)
func (app *serviceApplication) HandleFunc(methodString string, handler http.HandlerFunc) error {
	if app.httpServer == nil {
		return errors.New("http-сервер должен быть сначала инициализрован")
	}
	if methodString == "" {
		return errors.New("путь к методу должен быть определен")
	}
	if handler == nil {
		return errors.New("метод обработки должен быть определен")
	}
	app.httpServer.Handler.(*http.ServeMux).HandleFunc(methodString, handler)
	return nil
}

// Do - Метод запуска приложения
func (app *serviceApplication) Do() {
	log.Println("Application stated")
	closer.Bind(app.Close)

	// Запускаем http сервер если он определен
	if app.httpServer != nil {
		log.Printf("Starting http server at %s", app.httpServer.Addr)
		go func() {
			if err := app.httpServer.ListenAndServe(); err != nil {
				log.Fatal(err)
			}
		}()
	}

	closer.Hold()
}

// Close - метод закрытия приложения
func (app *serviceApplication) Close() {
	log.Println("Starting closing process")

	//0. Завершаем работу http сервера
	if app.httpServer != nil {
		log.Println("Shutdown http server")
		if ctx, err := context.WithTimeout(context.Background(), 30*time.Second); err == nil {
			app.httpServer.Shutdown(ctx)
		}
	}

	//1. Закрываем соединения для менеджеров очередей
	for id, manager := range app.messageManagers {
		log.Printf("Shutdown queue connection: %s", id)
		manager.Close()
	}
	//2. Закрываем соединения для менеджеров баз данных

	log.Println("Application is close")
}

var App serviceApplication

func init() {
	App.messageManagers = make(map[string]msg.MessageManager)
}
