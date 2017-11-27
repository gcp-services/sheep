package stats

import (
	"os"
	"sync"

	"github.com/denisbrodbeck/machineid"
	"github.com/rs/zerolog/log"
)

type Stats struct {
	hostname string
	uuid     string
	metrics  *sync.Map
}

func New() *Stats {
	host, err := os.Hostname()
	if err != nil {
		log.Panic().Err(err).Msg("error when obtaining hostname")
	}

	id, err := machineid.ProtectedID("sheep")
	if err != nil {
		log.Panic().Err(err).Msg("unable to obtain machine unique id")
	}

	s := &Stats{
		hostname: host,
		uuid:     id,
		metrics:  new(sync.Map),
	}
	return s
}

func (s *Stats) Gauge(name string, value float64) {
	s.metrics.Store(name, value)
}

func (s *Stats) Incr(name string, value float64) {
	v, _ := s.metrics.LoadOrStore(name, 0)
	s.metrics.Store(name, v.(float64)+value)
}

func (s *Stats) Get(name string) float64 {
	value, _ := s.metrics.LoadOrStore(name, 0)
	// TODO check type when adding different types
	return value.(float64)
}

// Save stats to the database
func (s *Stats) Save() {

}
