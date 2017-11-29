package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

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
		steps = append(steps, random(-10000, 10000))
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
			for cn := step; cn < 0; cn++ {
				opch <- -1
			}
		} else if step > 0 {
			for cn := step; cn > 0; cn-- {
				opch <- 1
			}
		}
	}
	log.Info().Msg("Okay, kicking off the operations!")

	ctx, cancel := context.WithCancel(context.Background())
	for tc := 0; tc <= 1; tc++ {
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

func work(ctx context.Context, opcn chan int) {
	client := &http.Client{}
	for {
		select {
		case op := <-opcn:
			if op != 0 {
				// placeholder
			}
			request, err := http.NewRequest("PUT", "http://localhost:5309/v1/incr", strings.NewReader("test"))
			if err != nil {
				panic(err)
			}
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
