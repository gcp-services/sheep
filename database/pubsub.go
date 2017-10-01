package database

import (
	"context"

	"cloud.google.com/go/pubsub"
)

type Pubsub struct {
	client *pubsub.Client
}

// SetupPubsub global client
func NewPubsub(project string) (*Pubsub, error) {
	client, err := pubsub.NewClient(context.Background(), project)

	if err != nil {
		return nil, err
	}
	return &Pubsub{
		client: client,
	}, nil
}

func (p *Pubsub) Read() {

}

func (p *Pubsub) Save(message *Message) error {
	return nil
}
