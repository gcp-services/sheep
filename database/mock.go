package database

import (
	"context"
	"errors"
	"sync"

	"cloud.google.com/go/spanner"
	"google.golang.org/grpc/codes"
)

type MockDatabase struct {
	db     map[string]int64
	exists map[string]bool
	log    map[string]bool
	lock   sync.Mutex
}

type MockQueue struct {
	queue []*Message
	c     chan bool
}

func NewMockDatabase() (Database, error) {
	return &MockDatabase{
		db:     make(map[string]int64),
		exists: make(map[string]bool),
		log:    make(map[string]bool),
		lock:   sync.Mutex{},
	}, nil
}

func NewMockQueue() (Stream, error) {
	return &MockQueue{
		c: make(chan bool),
	}, nil
}

func (db *MockDatabase) Save(m *Message) error {
	db.lock.Lock()
	defer db.lock.Unlock()
	if db.log[m.UUID] {
		return nil
	}
	key := m.Keyspace + m.Key + m.Name
	db.log[m.UUID] = true
	switch m.Operation {
	case "INCR":
		db.db[key]++
		db.exists[key] = true
	case "DECR":
		db.db[key]--
		db.exists[key] = true
	case "SET":
		db.db[key] = m.Value
		db.exists[key] = true
	default:
		return errors.New("invalid op")
	}
	return nil
}

func (db *MockDatabase) Read(m *Message) error {
	db.lock.Lock()
	defer db.lock.Unlock()
	if !db.exists[m.Keyspace+m.Key+m.Name] {
		return &spanner.Error{Code: codes.NotFound}
	}
	m.Value = db.db[m.Keyspace+m.Key+m.Name]
	return nil
}

func (q *MockQueue) Save(m *Message) error {
	q.queue = append(q.queue, m)
	q.c <- true
	return nil
}

func (q *MockQueue) Read(ctx context.Context, fn MessageFn) error {
	select {
	case <-q.c:
		go func() {
			var m *Message
			m, q.queue = q.queue[0], q.queue[1:]
			ok := fn(m)
			if !ok {
				q.Save(m)
			}
		}()
	}
	return nil
}

// TODO: implement cancel channel
func (q *MockQueue) StartWork(db Database) {
	go q.Read(context.Background(), func(msg *Message) bool {
		err := db.Save(msg)
		if err != nil {
			return false
		}
		return true
	})
}

func (q *MockQueue) StopWork() {

}
