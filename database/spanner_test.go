package database

import (
	"testing"
)

func TestSetupSpanner(t *testing.T) {
	// Remember for this test, these env vars must be set:
	// SHEEP_PROJECT
	// SHEEP_INSTANCE
	// SHEEP_DATABASE
	// TODO: Mock spanner :(
	//assert.Nil(t, SetupSpanner())
}
