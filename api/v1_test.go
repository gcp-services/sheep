package api

import (
	"bytes"
	"encoding/json"
	"net/http"
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

	q.StartWork(db)

	return New(&q, &db), nil
}

func setupRequest(m *database.Message, web *Handler) (*httptest.ResponseRecorder, echo.Context) {
	data, _ := json.Marshal(m)

	e := echo.New()
	req := httptest.NewRequest(echo.PUT, "/"+web.Version+"/"+m.Operation, bytes.NewBuffer(data))
	rec := httptest.NewRecorder()
	req.Header.Set("Content-Type", "application/json")

	c := e.NewContext(req, rec)
	c.SetPath("/" + web.Version + "/" + m.Operation)
	return rec, c
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

// Test Submit for expected results (but not incr/decr/set it self)
func TestSubmit(t *testing.T) {
	m := &database.Message{}

	web, err := setupWeb()
	assert.Nil(t, err)

	// Keyspace Check
	rec, c := setupRequest(m, web)
	if assert.NoError(t, web.Submit(c, "incr")) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, string(rec.Body.Bytes()), "missing keyspace field")
	}

	// Key check
	m.Keyspace = "testKeyspace"
	rec, c = setupRequest(m, web)
	if assert.NoError(t, web.Submit(c, "incr")) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, string(rec.Body.Bytes()), "missing key field")
	}

	// Name check
	m.Key = "testKey"
	rec, c = setupRequest(m, web)
	if assert.NoError(t, web.Submit(c, "incr")) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, string(rec.Body.Bytes()), "missing name field")
	}

	// UUID check
	m.Name = "testName"
	rec, c = setupRequest(m, web)
	if assert.NoError(t, web.Submit(c, "incr")) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, string(rec.Body.Bytes()), "missing uuid field")
	}

	// Good request
	m.UUID = "testUUID"
	rec, c = setupRequest(m, web)
	if assert.NoError(t, web.Submit(c, "incr")) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

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
