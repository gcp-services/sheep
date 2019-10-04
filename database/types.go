package database

import (
	"context"
)

type Stream interface {
	Save(*Message) error
	Read(context.Context, MessageFn) error
	StartWork(Database)
	StopWork()
}

type Database interface {
	Save(*Message) error
	Read(*Message) error
}

// Message struct for doing an operation.
type Message struct {
	UUID      string
	Keyspace  string
	Key       string
	Name      string
	Operation string
	Value     int64
}

type MessageFn func(*Message) bool

type ContextKey string
