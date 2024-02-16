package msg_srv

//MessageManager - интерфейс для реализации менеджера очередей
type MessageManager interface {
	ID() string
	CreateProducer(producerId, queueName string) error
	GetProducer(producerId string) (MessageProducer, bool)
	Consume(queueName string, handler MessageHandler) error
	Close()
}

//MessageProducer - интерфейс для реализации публикации сообщений в очередь
type MessageProducer interface {
	Manager() MessageManager
	SendMsg(key, value []byte, headers map[string]interface{}) error
}

//MessageHandler - тип для пользовательской функции обработки сообщения
type MessageHandler func(key, value []byte, headers map[string]interface{})
