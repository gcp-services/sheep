package database

import (
	"testing"

	"github.com/Cidan/sheep/config"
)

func TestSetupRabbitMQ(t *testing.T) {
	config.Setup("../config/")
	SetupRabbitMQ()
}
