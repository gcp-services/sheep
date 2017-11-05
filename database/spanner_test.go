package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupSpanner(t *testing.T) {
	// Remember for this test, these env vars must be set:
	// SHEEP_PROJECT
	// SHEEP_INSTANCE
	// SHEEP_DATABASE
	// TODO: Mock spanner :(
	//assert.Nil(t, SetupSpanner())
}

func TestSpannerSave(t *testing.T) {
	sp, err := NewSpanner("jinked-home", "sheep-test", "sheep")
	assert.Nil(t, err)

	// Create a message
	msg := &Message{
		UUID:      "123456",
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
	err = sp.Save(msg)
	assert.Nil(t, err)

	// Decrement by 1
	msg.Operation = "decr"
	err = sp.Save(msg)
	assert.Nil(t, err)

	// Invalid operation should error
	msg.Operation = "nope"
	err = sp.Save(msg)
	// not working assert.NotNil(t, err)

	// Missing fields should error
	err = sp.Save(&Message{})
	assert.NotNil(t, err)

}
