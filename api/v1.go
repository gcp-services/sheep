package api

import (
	"cloud.google.com/go/spanner"
	"github.com/Cidan/sheep/database"
	"github.com/labstack/echo"
	"google.golang.org/grpc/codes"
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

	// Operations
	e.PUT(h.path("/incr"), func(c echo.Context) error {
		return h.Submit(c, "incr")
	})
	e.PUT(h.path("/decr"), func(c echo.Context) error {
		return h.Submit(c, "decr")
	})
	e.PUT(h.path("/set"), func(c echo.Context) error {
		return h.Submit(c, "set")
	})
	return nil
}

func (h *Handler) Get(c echo.Context) error {
	msg := &database.Message{
		Keyspace: c.QueryParam("keyspace"),
		Key:      c.QueryParam("key"),
		Name:     c.QueryParam("name"),
	}
	err := h.Database.Read(msg)
	if err != nil {
		if spanner.ErrCode(err) == codes.NotFound {
			return c.JSON(404, echo.ErrNotFound)
		}
		return err
	}
	return c.JSON(200, msg)
}

func (h *Handler) Submit(c echo.Context, op string) error {
	msg := &database.Message{}
	if err := c.Bind(msg); err != nil {
		return err
	}

	if err := validateMessage(msg); err != nil {
		return c.JSON(400, err)
	}

	msg.Operation = op
	return h.Stream.Save(msg)
}

func validateMessage(msg *database.Message) error {
	if msg.Key == "" {
		return echo.NewHTTPError(400, "invalid payload, missing key field")
	}
	if msg.Keyspace == "" {
		return echo.NewHTTPError(400, "invalid payload, missing keyspace field")
	}
	if msg.Name == "" {
		return echo.NewHTTPError(400, "invalid payload, missing name field")
	}
	if msg.UUID == "" {
		return echo.NewHTTPError(400, "invalid payload, missing uuid field")
	}
	return nil
}
