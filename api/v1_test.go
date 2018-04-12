package api

import (
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Cidan/sheep/database"
	"github.com/labstack/echo"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func setupWeb() (*Handler, error) {
	db, err := database.NewMockDatabase()
	if err != nil {
		return nil, err
	}

	q, err := database.NewMockQueue()
	if err != nil {
		return nil, err
	}

	return New(&q, &db), nil
}

func TestNew(t *testing.T) {
	web, err := setupWeb()
	assert.Nil(t, err)
	assert.NotNil(t, web)
}

func TestRegister(t *testing.T) {
	web, err := setupWeb()
	assert.Nil(t, err)

	e := echo.New()
	e.Logger.SetOutput(log.Logger)
	e.HideBanner = true

	web.Register(e)
}

// TODO: Submit

func TestGet(t *testing.T) {
	web, err := setupWeb()
	assert.Nil(t, err)

	q := make(url.Values)
	q.Set("keyspace", "testKeyspace")
	q.Set("key", "testKey")
	q.Set("name", "testName")

	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/"+web.Version+"/get?"+q.Encode(), nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetPath("/" + web.Version + "/get")

	// Happy path, all is good
	assert.NoError(t, web.Get(c))
}
