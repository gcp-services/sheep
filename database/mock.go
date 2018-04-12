package database

import "context"

type MockDatabase struct {
}

type MockQueue struct {
}

func SetupMockDatabase() Database {
	return &MockDatabase{}
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
	return nil
}
func (db *MockDatabase) Read(m *Message) error {
	return nil
}

func (q *MockQueue) Save(m *Message) error {
	return nil
}

func (q *MockQueue) Read(ctx context.Context, fn MessageFn) error {
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
