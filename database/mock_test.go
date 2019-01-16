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

func TestMockDatabaseSaveInvalidOp(t *testing.T) {
	db, err := NewMockDatabase(false)
	assert.Nil(t, err)
	assert.NotNil(t, db)
	err = db.Save(&Message{})
	assert.EqualError(t, err, "invalid op")
}
