package rabbitmq_test

import (
	mq "async-arch/internal/lib/messages/rabbitmq"
	"async-arch/internal/sysenv"
	"testing"
)

var (
	result chan string
)

func CatchMessage(key, value []byte, headers map[string]interface{}) {
	result <- string(key) + "=" + string(value)
}

func TestRabbitMQ(t *testing.T) {
	result = make(chan string)

	serverAddr := sysenv.GetEnvValue("RABBITMQ_SERVER", "192.168.1.99")
	vhostName := sysenv.GetEnvValue("RABBITMQ_VHOST", "async_arch")

	manager, err := mq.CreateRabbitMQManager("asyncarch", "password", serverAddr, vhostName, mq.DEFAULT_PORT)
	if err != nil {
		t.Fatal(err)
	}

	err = manager.CreateProducer("test", "test")
	if err != nil {
		t.Fatal(err)
	}

	err = manager.Consume("test", CatchMessage)
	if err != nil {
		t.Fatal(err)
	}

	producer, _ := manager.GetProducer("test")

	err = producer.SendMsg([]byte("key"), []byte("value"), nil)
	if err != nil {
		t.Fatal(err)
	}

	msg := <-result
	if msg != "key=value" {
		t.Fatal("сообщение не получено")
	}

	manager.Close()
}
