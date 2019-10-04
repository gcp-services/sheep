package database

import (
	"database/sql"
	"fmt"

	// Import for postgres

	_ "github.com/lib/pq"
)

type CockroachDB struct {
	client   *sql.DB
	host     string
	username string
	password string
	dbname   string
	sslmode  string
	port     int
}

func NewCockroachDB(host, username, password, dbname, sslmode string, port int) (*CockroachDB, error) {
	c := &CockroachDB{
		host:     host,
		username: username,
		dbname:   dbname,
		sslmode:  sslmode,
		port:     port,
	}
	db, err := sql.Open("postgres", fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=%s port=%d",
		host, username, password, dbname, sslmode, port,
	))
	if err != nil {
		return nil, err
	}
	c.client = db

	return c, c.setupDatabase()
}

func (c *CockroachDB) setupDatabase() error {
	// Create our database and table if they don't exist.
	if _, err := c.client.Query(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", c.dbname)); err != nil {
		return err
	}
	return nil
}

func (c *CockroachDB) Read() {

}

func (c *CockroachDB) Save() {

}
