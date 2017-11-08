package database

import (
	"testing"

	"github.com/Cidan/sheep/config"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestSetupSpanner(t *testing.T) {
}

func TestSpannerSave(t *testing.T) {
	config.SetDefaults()
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

	// Set our counter to 0 and read it back
	err = sp.Save(msg)
	assert.Nil(t, err)
	err = sp.Read(msg)
	assert.Nil(t, err)
	assert.EqualValues(t, 0, msg.Value)

	// Increment by 1 and read it back
	msg.Operation = "incr"
	msg.UUID = uuid.NewV4().String()
	err = sp.Save(msg)
	assert.Nil(t, err)
	err = sp.Read(msg)
	assert.Nil(t, err)
	assert.EqualValues(t, 1, msg.Value)

	// Decrement by 1
	msg.Operation = "decr"
	msg.UUID = uuid.NewV4().String()
	err = sp.Save(msg)
	assert.Nil(t, err)
	err = sp.Read(msg)
	assert.Nil(t, err)
	assert.EqualValues(t, 0, msg.Value)

	// Invalid operation should error
	msg.Operation = "nope"
	msg.UUID = uuid.NewV4().String()
	err = sp.Save(msg)
	assert.NotNil(t, err)

	// Missing fields should error
	err = sp.Save(&Message{})
	assert.NotNil(t, err)

}
