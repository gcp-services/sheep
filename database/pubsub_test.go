package database

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPubsub(t *testing.T) {
	os.Setenv("PUBSUB_EMULATOR_HOST", "localhost:8085")
	_, err := NewPubsub("test")
	assert.Nil(t, err)
}
