package database

type Stream interface {
	Save(*Message) error
	Read(MessageFn) error
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
	Ack       chan bool
}

type MessageFn func(*Message) bool

type contextKey string
