package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/satori/go.uuid"

	"github.com/Cidan/sheep/database"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func main() {
	rand.Seed(time.Now().Unix())
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	log.Info().Msg("Starting up demo work")
	var steps []int
	for i := 1; i <= 10; i++ {
		steps = append(steps, random(-10, 10))
	}
	log.Info().Msg("Calculating total operations and final counter number")
	var totalSteps int
	var total int
	for _, step := range steps {
		total += step
		if step < 0 {
			totalSteps -= step
		} else {
			totalSteps += step
		}
	}
	log.Info().
		Int("steps", totalSteps).
		Int("total", total).Msg("Operations and total count")

	opch := make(chan int, totalSteps)
	for _, step := range steps {
		if step < 0 {
			log.Info().Int("count", step).Msg("Decrement by count")
			for cn := step; cn < 0; cn++ {
				opch <- -1
			}
		} else if step > 0 {
			log.Info().Int("count", step).Msg("Increment by count")
			for cn := step; cn > 0; cn-- {
				opch <- 1
			}
		}
	}
	log.Info().Msg("Okay, kicking off the operations!")
	reset()
	ctx, cancel := context.WithCancel(context.Background())
	for tc := 0; tc <= 4; tc++ {
		go work(ctx, opch)
	}
	for {
		time.Sleep(time.Second * 2)
		if len(opch) > 0 {
			fmt.Print(".")
		} else {
			break
		}
	}
	cancel()
}

func reset() {
	client := &http.Client{}
	msg := database.Message{
		UUID:     uuid.NewV4().String(),
		Keyspace: "test",
		Key:      "test",
		Name:     "test",
		Value:    0,
	}
	b, _ := json.Marshal(msg)

	request, _ := http.NewRequest("PUT", "http://localhost:5309/v1/set?direct=true", strings.NewReader(string(b)))
	request.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(request)
	log.Info().Int("status code", resp.StatusCode).Msg("reset counter to 0")
	if err != nil {
		panic(err)
	}
}

func work(ctx context.Context, opcn chan int) {
	msg := database.Message{
		Keyspace: "test",
		Key:      "test",
		Name:     "test",
	}
	client := &http.Client{}
	for {
		select {
		case op := <-opcn:
			if op != 0 {
				// placeholder
			}
			msg.UUID = uuid.NewV4().String()
			b, _ := json.Marshal(msg)
			var request *http.Request
			if op == -1 {
				request, _ = http.NewRequest("PUT", "http://localhost:5309/v1/decr", strings.NewReader(string(b)))
			} else {
				request, _ = http.NewRequest("PUT", "http://localhost:5309/v1/incr", strings.NewReader(string(b)))
			}
			request.Header.Add("Content-Type", "application/json")
			response, err := client.Do(request)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 200 {
				// retry
			}

		case <-ctx.Done():
			break
		}
	}
}
