package main

import (
	msg "async-arch/lib/msgconnect/rabbitmq"
)

func main() {
	manager := msg.CreateRabbitMQManager("asyncarch", "password", "192.168.1.99")

}
