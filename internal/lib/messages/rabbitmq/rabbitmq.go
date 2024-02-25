package rabbitmq

import (
	msg "async-arch/internal/lib/messages"
	"context"
	"errors"
	"fmt"
	"time"

	uuid "github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

const DEFAULT_PORT uint16 = 5672

// RabitMQManage - менеджер подключения к серверу RabbitMQ
type RabbitMQManager struct {
	id        string
	conn      *amqp.Connection
	channel   *amqp.Channel
	producers map[string]RabbitMessageProducer
	consumers map[string]RabbitMessageConsumer
}

// CreateRabbitMQManager - функция подключения к менеджеру очередей RabbitMQ
func CreateRabbitMQManager(user, password, server, vhost string, port uint16) (*RabbitMQManager, error) {
	//1. Создаем подключение к RabbitMQ
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/%s", user, password, server, port, vhost))
	if err != nil {
		return nil, err
	}

	//2. Создаем канал
	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQManager{
		id:        uuid.NewString(),
		conn:      conn,
		channel:   channel,
		producers: make(map[string]RabbitMessageProducer),
		consumers: make(map[string]RabbitMessageConsumer),
	}, nil
}

// CreateProducer - метод создания простого публкатора сообщений
func (man *RabbitMQManager) CreateProducer(producerId, queueName string) error {
	if _, ok := man.producers[producerId]; ok {
		return errors.New("публикатор с таким ID уже зарегистрирован")
	}

	queue, err := man.channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	producer := RabbitMessageProducer{
		manager: man,
		queue:   queue,
	}

	man.producers[producerId] = producer

	return nil
}

// GetProducer - метод возвращает публикатора по имени
func (man *RabbitMQManager) GetProducer(producerId string) (msg.MessageProducer, bool) {
	if producer, ok := man.producers[producerId]; ok {
		return &producer, ok
	} else {
		return nil, ok
	}
}

// Consume - метод создания простого читателя сообщений
func (man *RabbitMQManager) Consume(queueName string, handler msg.MessageHandler) error {
	id := uuid.NewString()
	queue, err := man.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	messages, err := man.channel.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	consumer := RabbitMessageConsumer{
		id:      id,
		manager: man,
		queue:   queue,
		cancel:  make(chan bool),
	}

	man.consumers[id] = consumer

	// Запускаем функцию проверки и чтения очереди
	go func() {

		for {
			select {
			// Читаем сообщения и запускаем их в обработку
			case message := <-messages:
				{
					go handler([]byte(message.MessageId), message.Body, message.Headers)
				}
			// Ждем сообщения о завершении и в случае чего завершаем выполнение процедуры
			case <-consumer.cancel:
				{
					return
				}
			}
		}
	}()
	return nil
}

// ID - метод возвращает уникальный идентификатор менеджера очередей
func (man *RabbitMQManager) ID() string {
	return man.id
}

// Close - метод закрытия менеджеров очередей
func (man *RabbitMQManager) Close() {
	for _, consumer := range man.consumers {
		consumer.Break()
	}
	man.channel.Close()
	man.conn.Close()
}

// RabbitMessageProducer - релизация публикатора для RabbitMQ (простой механизм)
type RabbitMessageProducer struct {
	manager *RabbitMQManager
	queue   amqp.Queue
}

// Manager - метод возвращает используемый менеджер очередей для публикатора
func (obj *RabbitMessageProducer) Manager() msg.MessageManager {
	return obj.manager
}

// SendMsg - метод публикации сообщения в очередь
func (obj *RabbitMessageProducer) SendMsg(key, value []byte, headers map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := obj.manager.channel.PublishWithContext(ctx, "", obj.queue.Name, false, false, amqp.Publishing{ContentType: "application/octet-stream", MessageId: string(key), Body: value, Headers: headers}); err != nil {
		return err
	} else {
		return nil
	}
}

// RabbitMessageConsumer - реализация читателя для RabbitMQ (простой механиз)
type RabbitMessageConsumer struct {
	id      string
	manager *RabbitMQManager
	queue   amqp.Queue
	cancel  chan bool
}

// Manager - метод возвращает используемый менеджер очередей для читателя
func (obj *RabbitMessageConsumer) Manager() msg.MessageManager {
	return obj.manager
}

// Close - метод остановки процесса чтения
func (obj *RabbitMessageConsumer) Break() {
	obj.cancel <- true
}
