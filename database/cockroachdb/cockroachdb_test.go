package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCockroachdb(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	_, err := NewCockroachDB("localhost", "root", "", "sheep", "disable", 26257)
	assert.Nil(t, err)
}
