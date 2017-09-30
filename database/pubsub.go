package database

import (
	"context"
	"os"

	"cloud.google.com/go/pubsub"
)

type Pubsub struct {
	client *pubsub.Client
}

// SetupPubsub global client
func NewPubsub() (*Pubsub, error) {
	client, err := pubsub.NewClient(context.Background(),
		os.Getenv("SHEEP_PROJECT"))

	if err != nil {
		return nil, err
	}
	return &Pubsub{
		client: client,
	}, nil
}

func (p *Pubsub) Read() {

}

func (p *Pubsub) Save() {

}
