package msg_srv

//MessageHandler - тип для пользовательской функции обработки сообщения
type MessageHandler func(key, value []byte, headers []string)

//MessageProducer - интерфейс для реализации публикации сообщений в очередь
type MessageProducer interface {
	SendMsg(key, value []byte, headers []string) error
}

//MessageConsumer - интерфейс для реализации чтения сообщений из очереди
type MessageConsumer interface {
	ReadMsg(message_handler MessageHandler) error
	Close()
}
