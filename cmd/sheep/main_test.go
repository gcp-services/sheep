package main

import (
	"testing"

	"github.com/Cidan/sheep/config"
	"github.com/Cidan/sheep/database"
	"github.com/Cidan/sheep/util"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestSetupLogging(t *testing.T) {
	setupLogging()
}

func TestStartGrpc(t *testing.T) {
	config.Setup("")
	db, err := database.NewMockDatabase(false)
	assert.Nil(t, err)

	stream, err := database.NewMockQueue(false)
	assert.Nil(t, err)

	go startGrpc(stream, db)
	assert.True(t, util.WaitForPort("localhost", viper.GetInt("service.port"), 6))

	go startWeb(stream, db)
	assert.True(t, util.WaitForPort("localhost", viper.GetInt("service.rest"), 6))
}

func TestSetupDatabase(t *testing.T) {
	config.Setup("")
	_, err := database.NewMockDatabase(false)
	assert.Nil(t, err)
}

func TestSetupQueue(t *testing.T) {
	config.Setup("")
	_, err := database.NewMockQueue(false)
	assert.Nil(t, err)
}
