package database

import (
	"testing"

	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestSetupSpanner(t *testing.T) {
}

func TestSpannerSave(t *testing.T) {
	sp, err := NewSpanner("jinked-home", "sheep-test", "sheep")
	assert.Nil(t, err)

	// Create a message
	msg := &Message{
		UUID:      uuid.NewV4().String(),
		Keyspace:  "test",
		Key:       "test",
		Name:      "some counter",
		Operation: "set",
		Value:     0,
	}

	// Set our counter to 0
	err = sp.Save(msg)
	assert.Nil(t, err)

	// TODO: READ

	// Increment by 1
	msg.Operation = "incr"
	msg.UUID = uuid.NewV4().String()
	err = sp.Save(msg)
	assert.Nil(t, err)

	// Decrement by 1
	msg.Operation = "decr"
	msg.UUID = uuid.NewV4().String()
	err = sp.Save(msg)
	assert.Nil(t, err)

	// Invalid operation should error
	msg.Operation = "nope"
	msg.UUID = uuid.NewV4().String()
	err = sp.Save(msg)
	assert.NotNil(t, err)

	// Missing fields should error
	err = sp.Save(&Message{})
	assert.NotNil(t, err)

}
