package spanner

import (
	"testing"

	"github.com/Cidan/sheep/config"
	"github.com/Cidan/sheep/database"
	"github.com/gcpug/handy-spanner/fake"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/option"
)

func SetupFake() {

}
func TestSetupSpanner(t *testing.T) {
}

func TestSpannerSave(t *testing.T) {

	srv, conn, err := fake.Run()
	assert.Nil(t, err)
	assert.NotNil(t, srv)

	config.SetDefaults()
	sp, err := New("jinked-home", "sheep-test", "sheep", option.WithGRPCConn(conn))
	assert.Nil(t, err)

	// Create a message
	msg := &database.Message{
		UUID:      uuid.NewV4().String(),
		Keyspace:  "test",
		Key:       "test",
		Name:      "some counter",
		Operation: "SET",
		Value:     0,
	}

	// Set our counter to 0 and read it back
	err = sp.Save(msg)
	assert.Nil(t, err)
	err = sp.Read(msg)
	assert.Nil(t, err)
	assert.EqualValues(t, 0, msg.Value)

	// Increment by 1 and read it back
	msg.Operation = "INCR"
	msg.UUID = uuid.NewV4().String()
	err = sp.Save(msg)
	assert.Nil(t, err)
	err = sp.Read(msg)
	assert.Nil(t, err)
	assert.EqualValues(t, 1, msg.Value)

	// Decrement by 1
	msg.Operation = "DECR"
	msg.UUID = uuid.NewV4().String()
	err = sp.Save(msg)
	assert.Nil(t, err)
	err = sp.Read(msg)
	assert.Nil(t, err)
	assert.EqualValues(t, 0, msg.Value)

	// Invalid operation should error
	msg.Operation = "nope"
	msg.UUID = uuid.NewV4().String()
	err = sp.Save(msg)
	assert.NotNil(t, err)

	// Missing fields should error
	err = sp.Save(&database.Message{})
	assert.NotNil(t, err)

}
