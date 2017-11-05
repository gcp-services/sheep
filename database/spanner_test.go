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
	err = sp.Save(&Message{
		UUID:      "1234",
		Keyspace:  "test",
		Key:       "test",
		Name:      "some counter",
		Operation: "incr",
	})
	assert.Nil(t, err)
}
