package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMockDatabase(t *testing.T) {
	db, err := NewMockDatabase(false)
	assert.Nil(t, err)
	assert.NotNil(t, db)
}

func TestNewMockQueue(t *testing.T) {
	q, err := NewMockQueue(false)
	assert.Nil(t, err)
	assert.NotNil(t, q)
}

func TestMockDatabaseSaveError(t *testing.T) {
	db, err := NewMockDatabase(true)
	assert.Nil(t, err)
	assert.NotNil(t, db)
	err = db.Save(&Message{})
	assert.Error(t, err)
}

func TestMockDatabaseSave(t *testing.T) {
	db, err := NewMockDatabase(false)
	msg := &Message{
		UUID:     "1",
		Keyspace: "test",
		Key:      "test",
		Name:     "test",
	}
	assert.Nil(t, err)
	assert.NotNil(t, db)
	err = db.Save(msg)
	assert.EqualError(t, err, "invalid op")

	// Test INCR
	msg.Operation = "INCR"
	msg.UUID = "2"
	err = db.Save(msg)
	assert.Nil(t, err)
	err = db.Read(msg)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), msg.Value)

	// Test duplicate op
	msg.Operation = "INCR"
	err = db.Save(msg)
	assert.Nil(t, err)
	err = db.Read(msg)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), msg.Value)

	// Test DECR
	msg.Operation = "DECR"
	msg.UUID = "3"
	err = db.Save(msg)
	assert.Nil(t, err)
	err = db.Read(msg)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), msg.Value)

	// Test SET
	msg.Operation = "SET"
	msg.UUID = "4"
	msg.Value = 10
	err = db.Save(msg)
	assert.Nil(t, err)
	err = db.Read(msg)
	assert.Nil(t, err)
	assert.Equal(t, int64(10), msg.Value)
}

func TestMockDatabaseReadError(t *testing.T) {
	db, err := NewMockDatabase(true)
	assert.Nil(t, err)
	assert.NotNil(t, db)
	err = db.Read(&Message{})
	assert.Error(t, err)
}

func TestMockDatabaseRead(t *testing.T) {
	db, err := NewMockDatabase(false)
	assert.Nil(t, err)
	assert.NotNil(t, db)
	msg := &Message{}
	// Test Read Not Found
	err = db.Read(msg)
	assert.Error(t, err)

	msg.Key = "test"
	msg.Keyspace = "test"
	msg.Name = "test"
	msg.Operation = "INCR"
	msg.UUID = "1"

	// Test Read
	err = db.Save(msg)
	assert.Nil(t, err)
	err = db.Read(msg)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), msg.Value)
}
