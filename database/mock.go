package database

import (
	"context"
	"sync"
)

type MockDatabase struct {
	db   map[string]int64
	log  map[string]bool
	lock sync.Mutex
}

type MockQueue struct {
	queue []*Message
	c     chan bool
}

func SetupMockDatabase() Database {
	return &MockDatabase{
		db:   make(map[string]int64),
		log:  make(map[string]bool),
		lock: sync.Mutex{},
	}
}

func SetupMockQueue() Stream {
	return &MockQueue{}
}

func NewMockDatabase() (Database, error) {
	return &MockDatabase{}, nil
}

func NewMockQueue() (Stream, error) {
	return &MockQueue{}, nil
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
	case "incr":
		db.db[key]++
	case "decr":
		db.db[key]--
	case "set":
		db.db[key] = m.Value
	}
	return nil
}

func (db *MockDatabase) Read(m *Message) error {
	db.lock.Lock()
	defer db.lock.Unlock()
	m.Value = db.db[m.Keyspace+m.Key+m.Name]
	return nil
}

func (q *MockQueue) Save(m *Message) error {
	q.queue = append(q.queue, m)
	q.c <- true
	return nil
}

func (q *MockQueue) Read(ctx context.Context, fn MessageFn) error {
	switch {
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
