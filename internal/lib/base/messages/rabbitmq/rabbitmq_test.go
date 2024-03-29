package rabbitmq_test

import (
	mq "async-arch/internal/lib/base/messages/rabbitmq"
	ou "async-arch/internal/lib/osutils"
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

	serverAddr := ou.GetEnvValue("RABBITMQ_SERVER", "localhost")
	vhostName := ou.GetEnvValue("RABBITMQ_VHOST", "async_arch")

	manager, err := mq.CreateRabbitMQManager("test", "asyncarch", "password", serverAddr, vhostName, mq.DEFAULT_PORT)
	if err != nil {
		t.Fatal(err)
	}

	err = manager.AddProducer("test", "test")
	if err != nil {
		t.Fatal(err)
	}

	err = manager.Consume("test", CatchMessage)
	if err != nil {
		t.Fatal(err)
	}

	producer, _ := manager.GetProducer("test")

	err = producer.ProduceMessage([]byte("key"), []byte("value"), nil)
	if err != nil {
		t.Fatal(err)
	}

	msg := <-result
	if msg != "key=value" {
		t.Fatal("сообщение не получено")
	}

	manager.Close()
}
