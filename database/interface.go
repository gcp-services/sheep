package database

type Stream interface {
	Save()
	Read()
}

type Database interface {
	Save()
	Read()
}
