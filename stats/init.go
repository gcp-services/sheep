package stats

import (
	"time"

	"github.com/rs/zerolog/log"
)

var metrics *Stats

func Setup() {
	log.Info().Msg("Setting up internal stats tracking")
	metrics = New()
	go func() {
		for {
			metrics.Save()
			time.Sleep(5 * time.Second)
		}
	}()
}

func Gauge(name string, value float64) {
	if metrics == nil {
		return
	}
	metrics.Gauge(name, value)
}

func Incr(name string, value float64) {
	if metrics == nil {
		return
	}
	metrics.Incr(name, value)
}

func Get(name string) float64 {
	if metrics == nil {
		return 0
	}
	return metrics.Get(name)
}
