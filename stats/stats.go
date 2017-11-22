package stats

import (
	"sync"

	"github.com/Cidan/sheep/database"
)

type Stats struct {
	db      database.Database
	metrics *sync.Map
}

func New(db database.Database) *Stats {
	s := &Stats{
		db:      db,
		metrics: new(sync.Map),
	}
	return s
}

func (s *Stats) Gauge(name string, value float64) {
	s.metrics.Store(name, value)
}

func (s *Stats) Get(name string) float64 {
	value, _ := s.metrics.LoadOrStore(name, 0)
	// TODO check type when adding different types
	return value.(float64)
}

// Save stats to the database
func (s *Stats) Save() {

}
