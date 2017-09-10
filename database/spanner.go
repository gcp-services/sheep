package database

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/spanner"
)

// Spanner is the global spanner client.
var Spanner *spanner.Client

// SetupSpanner initializes the spanner client.
func SetupSpanner() error {
	db := fmt.Sprintf("projects/%s/instances/%s/databases/%s",
		os.Getenv("SHEEP_PROJECT"),
		os.Getenv("SHEEP_INSTANCE"),
		os.Getenv("SHEEP_DATABASE"))

	client, err := spanner.NewClient(context.Background(), db)
	if err != nil {
		return err
	}
	Spanner = client
	return nil
}
