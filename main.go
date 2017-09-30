package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Cidan/sheep/api"
	"github.com/Cidan/sheep/config"
	"github.com/Cidan/sheep/database"
	"github.com/labstack/echo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
)

var e *echo.Echo

func main() {
	setupLogging()
	config.Setup("")
	stream, err := setupQueue()
	if err != nil {
		log.Panic().Err(err).Msg("Could not setup queue")
	}
	database, err := setupDatabase()
	if err != nil {
		log.Panic().Err(err).Msg("Could not setup database")
	}
	setupWebserver(stream, database)
}

func setupDatabase() (database.Database, error) {
	if viper.GetBool("cockroachdb.enabled") {
		return nil, nil
	}

	if viper.GetBool("spanner.enabled") {
		client, err := database.NewSpanner(
			viper.GetString("spanner.project"),
			viper.GetString("spanner.instance"),
			viper.GetString("spanner.database"),
		)
		if err != nil {
			log.Panic().Err(err).Msg("Could not start Spanner connection")
		}
		return client, err
	}

	return nil, fmt.Errorf("no database enabled")
}

func setupQueue() (database.Stream, error) {
	if viper.GetBool("rabbitmq.enabled") {
		return nil, nil
	}

	if viper.GetBool("pubsub.enabled") {
		pubsub, err := database.NewPubsub()
		if err != nil {
			log.Panic().Err(err).Msg("Could not start Pubsub connection")
		}
		return pubsub, err
	}

	return nil, fmt.Errorf("no queue setup")
}

func setupWebserver(stream database.Stream, database database.Database) {
	e = echo.New()
	e.Logger.SetOutput(log.Logger)

	e.HideBanner = true

	e.GET("/healthz", func(c echo.Context) error {
		return c.String(200, "ok")
	})

	// Create our v1 API
	v1 := api.New()
	v1.Register(e)
	log.Info().
		Int("port", 5309).
		Msg("Starting webserver")
	err := e.Start(":5309")
	if strings.Contains(err.Error(), "Server closed") {
		log.Info().
			Msg("Server has shutdown.")
		return
	}
	log.Fatal().Err(err)
}

func stopWebserver() {

}
func setupLogging() {
	// If we're in a terminal, pretty print
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
		log.Info().Msg("Detected terminal, pretty logging enabled.")
	}
}
