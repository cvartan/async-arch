package base_test

import (
	"async-arch/internal/lib/base"
	mq "async-arch/internal/lib/messages/rabbitmq"
	"async-arch/internal/sysenv"
	"testing"
)

func TestApp(t *testing.T) {
	defer base.App.Close()
	if manager, err := mq.CreateRabbitMQManager("asyncarch", "password", sysenv.GetEnvValue("RABBITMQ_SERVER", "192.168.1.99"), sysenv.GetEnvValue("RABBITMQ_VHOST", "async_arch"), mq.DEFAULT_PORT); err == nil {
		base.App.RegisterMessageManager(manager)
	} else {
		t.Fatal(err)
	}
}
