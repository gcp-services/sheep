package main

import (
	"os"

	"github.com/Cidan/sheep/api"
	"github.com/Cidan/sheep/database"
	"github.com/labstack/echo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	setupLogging()
	database.SetupSpanner()
	database.SetupPubsub()
	setupWebserver()
}

func setupWebserver() {
	e := echo.New()
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
	e.Logger.Fatal(e.Start(":5309"))
}

func setupLogging() {
	// If we're in a terminal, pretty print
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
		log.Info().Msg("Detected terminal, pretty logging enabled.")
	}
}
