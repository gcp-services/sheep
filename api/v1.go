package api

import (
	"cloud.google.com/go/spanner"
	"github.com/Cidan/sheep/database"
	"github.com/Cidan/sheep/stats"
	"github.com/labstack/echo"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
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
			stats.Incr("http.get.404", 1)
			return c.JSON(404, echo.ErrNotFound)
		}
		stats.Incr("http.get.500", 1)
		return err
	}
	stats.Incr("http.get.200", 1)
	return c.JSON(200, msg)
}

func (h *Handler) Submit(c echo.Context, op string) error {
	stat := "http." + op
	msg := &database.Message{}
	if err := c.Bind(msg); err != nil {
		return err
	}

	msg.Operation = op

	if err := validateMessage(msg); err != nil {
		log.Error().Err(err).Msg("unable to put message")
		stats.Incr(stat+".400", 1)
		return c.JSON(400, err)
	}

	if viper.GetBool("direct") || c.QueryParam("direct") == "true" {
		if err := h.Database.Save(msg); err != nil {
			stats.Incr(stat+".500", 1)
			return err
		}
		stats.Incr(stat+".200", 1)
		return nil
	}

	if err := h.Stream.Save(msg); err != nil {
		stats.Incr(stat+".500", 1)
		return err
	}
	stats.Incr(stat+".200", 1)
	return nil
}

func validateMessage(msg *database.Message) error {
	if msg.Keyspace == "" {
		return echo.NewHTTPError(400, "invalid payload, missing keyspace field")
	}
	if msg.Key == "" {
		return echo.NewHTTPError(400, "invalid payload, missing key field")
	}
	if msg.Name == "" {
		return echo.NewHTTPError(400, "invalid payload, missing name field")
	}
	if msg.UUID == "" {
		return echo.NewHTTPError(400, "invalid payload, missing uuid field")
	}
	return nil
}
