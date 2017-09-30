package api

import (
	"github.com/Cidan/sheep/database"
	"github.com/labstack/echo"
)

type Handler struct {
	Version  string
	Database database.Database
	Stream   database.Stream
}

func New(stream *database.Stream, db *database.Database) *Handler {
	return &Handler{
		Version:  "v1",
		Database: *db,
		Stream:   *stream,
	}
}

func (h *Handler) path(path string) string {
	return "/" + h.Version + path
}

func (h *Handler) Register(e *echo.Echo) error {
	e.GET(h.path("/get"), h.Get)
	e.PUT(h.path("/incr"), func(c echo.Context) error {
		return h.Submit(c, "incr")
	})
	e.PUT(h.path("/decr"), func(c echo.Context) error {
		return h.Submit(c, "decr")
	})
	return nil
}

func (h *Handler) Get(c echo.Context) error {
	return nil
}

func (h *Handler) Submit(c echo.Context, op string) error {
	msg := &database.Message{}
	c.Bind(msg)
	msg.Operation = op
	return h.Stream.Save(msg)
}
