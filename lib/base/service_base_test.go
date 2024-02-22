package base_test

import (
	"async-arch/lib/base"
	mq "async-arch/lib/messages/rabbitmq"
	"async-arch/util"
	"testing"
)

func TestApp(t *testing.T) {
	defer base.App.Close()
	if manager, err := mq.CreateRabbitMQManager("asyncarch", "password", util.GetEnvValue("RABBITMQ_SERVER", "192.168.1.99"), util.GetEnvValue("RABBITMQ_VHOST", "async_arch"), mq.DEFAULT_PORT); err == nil {
		base.App.RegisterMessageManager(manager)
	} else {
		t.Fatal(err)
	}
}
