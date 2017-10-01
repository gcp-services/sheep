package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupRabbitMQ(t *testing.T) {
	_, err := NewRabbitMQ([]string{"amqp://localhost", "amqp://localhost"})
	assert.Nil(t, err)
}
