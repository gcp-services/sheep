package database

import (
	"context"
	"os"

	"cloud.google.com/go/pubsub"
)

// Pubsub global client
var Pubsub *pubsub.Client

// SetupPubsub global client
func SetupPubsub() error {
	client, err := pubsub.NewClient(context.Background(),
		os.Getenv("SHEEP_PROJECT"))

	if err != nil {
		return err
	}

	Pubsub = client
	return nil
}
