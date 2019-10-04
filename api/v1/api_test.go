package v1

import (
	"context"
	"testing"

	"github.com/Cidan/sheep/database/mock"
	"github.com/stretchr/testify/assert"
)

func TestGetWithNotFound(t *testing.T) {
	in := &Counter{
		Keyspace: "test",
		Key:      "test",
		Name:     "test",
	}

	db, err := mock.NewMockDatabase(false)
	assert.Nil(t, err)
	api := &API{
		Database: db,
	}

	res, err := api.Get(context.Background(), in)
	assert.NotNil(t, err)
	assert.Zero(t, res.GetValue())
}

func TestGetDirect(t *testing.T) {
	in := &Counter{
		Keyspace:  "test",
		Key:       "test",
		Name:      "test",
		Uuid:      "123",
		Operation: Counter_INCR,
		Direct:    true,
	}
	db, err := mock.NewMockDatabase(false)
	assert.Nil(t, err)
	api := &API{
		Database: db,
	}
	res, err := api.Update(context.Background(), in)
	assert.Nil(t, err)
	assert.Equal(t, "", res.GetError())

	res, err = api.Get(context.Background(), in)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), res.GetValue())
}

func TestGetStream(t *testing.T) {
	in := &Counter{
		Keyspace:  "test",
		Key:       "test",
		Name:      "test",
		Uuid:      "123",
		Operation: Counter_INCR,
		Direct:    false,
	}

	db, err := mock.NewMockDatabase(false)
	assert.Nil(t, err)

	stream, err := mock.NewMockQueue(false)
	assert.Nil(t, err)

	stream.StartWork(db)

	api := &API{
		Database: db,
		Stream:   stream,
	}

	res, err := api.Update(context.Background(), in)
	assert.Nil(t, err)
	assert.Equal(t, "", res.GetError())

	res, err = api.Get(context.Background(), in)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), res.GetValue())
}

func TestGetError(t *testing.T) {
	in := &Counter{
		Keyspace:  "test",
		Key:       "test",
		Name:      "test",
		Uuid:      "123",
		Operation: Counter_INCR,
		Direct:    false,
	}

	db, err := mock.NewMockDatabase(true)
	assert.Nil(t, err)

	stream, err := mock.NewMockQueue(false)
	assert.Nil(t, err)

	stream.StartWork(db)
	api := &API{
		Database: db,
		Stream:   stream,
	}
	res, err := api.Update(context.Background(), in)
	assert.Nil(t, err)
	assert.Equal(t, "", res.GetError())

	res, err = api.Get(context.Background(), in)
	assert.NotNil(t, err)
	assert.Equal(t, int64(0), res.GetValue())
}

func TestUpdateError(t *testing.T) {
	in := &Counter{
		Keyspace:  "test",
		Key:       "test",
		Name:      "test",
		Uuid:      "123",
		Operation: Counter_INCR,
		Direct:    true,
	}

	db, err := mock.NewMockDatabase(true)
	assert.Nil(t, err)

	api := &API{
		Database: db,
	}

	res, err := api.Update(context.Background(), in)
	assert.NotNil(t, err)
	assert.Equal(t, "", res.GetError())
}
