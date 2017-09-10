package api

import "github.com/labstack/echo"

type Handler struct {
	Version string
}

func New() *Handler {
	return &Handler{
		Version: "v1",
	}
}

func (h *Handler) path(path string) string {
	return "/" + h.Version + path
}

func (h *Handler) Register(e *echo.Echo) error {
	e.GET(h.path("/get"), h.Get)
	return nil
}

func (h *Handler) Get(c echo.Context) error {
	return nil
}
