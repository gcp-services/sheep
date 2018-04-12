package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var rmq *RabbitMQ

func TestSetupRabbitMQ(t *testing.T) {
	if !acc {
		return
	}
	r, err := NewRabbitMQ([]string{"amqp://localhost", "amqp://localhost"})
	assert.Nil(t, err)
	rmq = r
}

func TestSave(t *testing.T) {
	if !acc {
		return
	}
	msg := &Message{}
	rmq.Save(msg)

}
