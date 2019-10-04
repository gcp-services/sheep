package rabbitmq

import (
	"testing"

	"github.com/Cidan/sheep/database"
	"github.com/stretchr/testify/assert"
)

var rmq *RabbitMQ

func TestSetupRabbitMQ(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	r, err := NewRabbitMQ([]string{"amqp://localhost", "amqp://localhost"})
	assert.Nil(t, err)
	rmq = r
}

func TestSave(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	msg := &database.Message{}
	rmq.Save(msg)

}
