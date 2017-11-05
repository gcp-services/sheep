package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPubsub(t *testing.T) {
	// The pubsub emulator doesn't support exists() checks it seems.
	//os.Setenv("PUBSUB_EMULATOR_HOST", "localhost:8085")
	_, err := NewPubsub("jinked-home", "sheep", "sheep")
	assert.Nil(t, err)
}
