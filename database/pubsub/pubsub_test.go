package pubsub

import (
	"context"
	"testing"

	"github.com/Cidan/sheep/database"
	"github.com/stretchr/testify/assert"
)

func TestNewPubsub(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	// The pubsub emulator doesn't support exists() checks it seems.
	//os.Setenv("PUBSUB_EMULATOR_HOST", "localhost:8085")
	_, err := New("jinked-home", "tests", "tests")
	assert.Nil(t, err)
}

func TestPubsubSaveAndRead(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	p, err := New("jinked-home", "tests", "tests")
	assert.Nil(t, err)

	err = p.Save(&database.Message{
		UUID:  "1234",
		Value: 1337,
	})
	assert.Nil(t, err)

	c := make(chan bool)
	go p.Read(context.Background(), func(msg *database.Message) bool {
		if msg.UUID == "1234" && msg.Value == 1337 {
			c <- true
			return true
		}
		c <- false
		return false
	})

	assert.True(t, <-c)
}
