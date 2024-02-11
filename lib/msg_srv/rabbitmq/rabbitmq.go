package rabbitmq

import (
	msg "async-arch/lib/msg_srv"
)

// RabbitMessageProducer - релизация публикатора для RabbitMQ (простой механизм)
type RabbitMessageProducer struct {
	queue string
}

// SendMsg - метод публикации сообщения в очередь
func (obj *RabbitMessageProducer) SendMsg(key, value []byte, headers []string) error {
	return nil
}

// RabbitMessageConsumer - реализация читателя для RabbitMQ (простой механиз)
type RabbitMessageConsumer struct {
	queue      string
	cancelFlag bool
}

// ReadMsg - метод старта процесса чтения сообщений из очереди.
func (obj *RabbitMessageConsumer) ReadMsg(handler msg.MessageHandler) error {
	obj.cancelFlag = false

	// Запускаем функцию проверки и чтения очереди
	go func() {
		if obj.cancelFlag {
			return
		}
	}()

	return nil
}

// Close - метод остановки процесса чтения
func (obj *RabbitMessageConsumer) Close() {
	obj.cancelFlag = true
}
