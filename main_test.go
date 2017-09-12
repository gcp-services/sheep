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
	go setupWebserver()
	assert.True(t, util.WaitForPort("localhost", 5309, 5))
	e.Shutdown(context.Background())
}

func TestSetupDatabase(t *testing.T) {
	config.Setup()
	assert.Nil(t, setupDatabase())
}

func TestSetupQueue(t *testing.T) {
	config.Setup()
	assert.Nil(t, setupQueue())
}
