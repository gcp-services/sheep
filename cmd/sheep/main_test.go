package main

import (
	"context"
	"testing"

	"github.com/Cidan/sheep/config"
	"github.com/Cidan/sheep/util"
	"github.com/stretchr/testify/assert"
)

func TestSetupLogging(t *testing.T) {
	setupLogging()
}
func TestSetupWebserver(t *testing.T) {
	config.Setup("")
	db, err := setupDatabase()
	assert.Nil(t, err)

	stream, err := setupQueue()
	assert.Nil(t, err)

	go setupWebserver(stream, db)
	assert.True(t, util.WaitForPort("localhost", 5309, 5))
	e.Shutdown(context.Background())
}

func TestSetupDatabase(t *testing.T) {
	config.Setup("")
	_, err := setupDatabase()
	assert.Nil(t, err)
}

func TestSetupQueue(t *testing.T) {
	config.Setup("")
	_, err := setupQueue()
	assert.Nil(t, err)
}
