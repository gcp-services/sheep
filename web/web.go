package web

import (
	"github.com/labstack/echo"
	"github.com/rs/zerolog/log"
)

type Handler struct {
}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Register(e *echo.Echo) error {
	log.Info().Msg("Registering UI")
	e.Static("/ui", "web/assets/dist")
	return nil
}

// TODO: function for serving files compiled into binary
