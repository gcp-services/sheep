package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var rmq *RabbitMQ

func TestSetupRabbitMQ(t *testing.T) {
	r, err := NewRabbitMQ([]string{"amqp://localhost", "amqp://localhost"})
	assert.Nil(t, err)
	rmq = r
}

func TestSave(t *testing.T) {
	msg := &Message{}
	rmq.Save(msg)

}
