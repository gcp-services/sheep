package main

import (
	"os"
	"strings"

	"github.com/Cidan/sheep/api"
	"github.com/Cidan/sheep/database"
	"github.com/labstack/echo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	setupLogging()

	err := database.SetupSpanner()
	if err != nil {
		log.Panic().Err(err).Msg("Could not start Spanner connection")
	}

	err = database.SetupPubsub()
	if err != nil {
		log.Panic().Err(err).Msg("Could not start Pubsub connection")
	}

	setupWebserver()
}

var e *echo.Echo

func setupWebserver() {
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
