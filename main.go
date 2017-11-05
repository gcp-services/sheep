package main

import (
	"errors"
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
		return nil, errors.New("CockroachDB support is not ready")
		/*
			return database.NewCockroachDB(
				viper.GetString("cockroachdb.host"),
				viper.GetString("cockroachdb.username"),
				viper.GetString("cockroachdb.password"),
				viper.GetString("cockroachdb.dbname"),
				viper.GetString("cockroachdb.sslmode"),
				viper.GetInt("cockroachdb.port"),
			)
		*/
	}

	if viper.GetBool("spanner.enabled") {
		log.Info().Msg("Setting up Spanner connection and schema")
		return database.NewSpanner(
			viper.GetString("spanner.project"),
			viper.GetString("spanner.instance"),
			viper.GetString("spanner.database"),
		)
	}

	return nil, fmt.Errorf("no database enabled")
}

func setupQueue() (database.Stream, error) {
	if viper.GetBool("rabbitmq.enabled") {
		return nil, errors.New("RabbitMQ support is not ready")
		/*
			return database.NewRabbitMQ(
				viper.GetStringSlice("rabbitmq.hosts"),
			)
		*/
	}

	if viper.GetBool("pubsub.enabled") {
		log.Info().Msg("Setting up Pub/Sub connection, topic, and subscription")
		return database.NewPubsub(
			viper.GetString("pubsub.project"),
			viper.GetString("pubsub.topic"),
			viper.GetString("pubsub.subscription"),
		)
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
	v1 := api.New(&stream, &database)
	v1.Register(e)
	log.Info().
		Int("port", 5309).
		Msg("Started webserver, you're good to go!")
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
