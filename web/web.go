package web

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/rakyll/statik/fs"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	_ "github.com/Cidan/sheep/statik"
)

type Handler struct {
}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Embed(e *echo.Echo) error {
	log.Info().Msg("Loading embedded UI")
	statikFS, err := fs.New()
	if err != nil {
		return err
	}
	assets := http.FileServer(statikFS)
	e.GET("/ui/*", echo.WrapHandler(http.StripPrefix("/ui/", assets)))
	return nil
}

func (h *Handler) Register(e *echo.Echo) error {
	path := viper.GetString("ui.path")
	log.Info().Str("path", path).Msg("Serving UI from path")
	e.Static("/ui", path)
	return nil
}
