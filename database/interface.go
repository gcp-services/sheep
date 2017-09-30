package database

type Stream interface {
	Save(*Message) error
	Read()
}

type Database interface {
	Save()
	Read()
}

// Message struct for doing an operation.
type Message struct {
	UUID      string
	Keyspace  string
	Key       string
	Name      string
	Operation string
}
