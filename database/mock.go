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

func NewMockDatabase() (*MockDatabase, error) {
	return &MockDatabase{}, nil
}

func NewMockQueue() (*MockQueue, error) {
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
