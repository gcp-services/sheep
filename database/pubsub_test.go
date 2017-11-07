package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPubsub(t *testing.T) {
	// The pubsub emulator doesn't support exists() checks it seems.
	//os.Setenv("PUBSUB_EMULATOR_HOST", "localhost:8085")
	_, err := NewPubsub("jinked-home", "tests", "tests")
	assert.Nil(t, err)
}

func TestPubsubSaveAndRead(t *testing.T) {
	p, err := NewPubsub("jinked-home", "tests", "tests")
	assert.Nil(t, err)

	err = p.Save(&Message{
		UUID:  "1234",
		Value: 1337,
	})
	assert.Nil(t, err)

	c := make(chan bool)
	go p.Read(func(msg *Message) bool {
		if msg.UUID == "1234" && msg.Value == 1337 {
			c <- true
			return true
		}
		c <- false
		return false
	})

	assert.True(t, <-c)
}
