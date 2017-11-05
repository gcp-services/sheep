package database

type Stream interface {
	Save(*Message) error
	Read() (chan *Message, error)
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
	Ack       chan bool
}
