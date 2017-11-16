package web

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/rakyll/statik/fs"
	"github.com/rs/zerolog/log"

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
	log.Info().Msg("Serving UI from web/assets/dist")
	e.Static("/ui", "web/assets/dist")
	return nil
}
