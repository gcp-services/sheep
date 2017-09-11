package database

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/spanner"
	database "cloud.google.com/go/spanner/admin/database/apiv1"
	adminpb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
)

// Spanner is the global spanner client.
var Spanner *spanner.Client

func createSpannerDatabase(ctx context.Context, admin *database.DatabaseAdminClient, project, instance, db string) error {
	// Create our database if it doesn't exist.
	_, err := admin.GetDatabase(ctx, &adminpb.GetDatabaseRequest{
		Name: "projects/" + project + "/instances/" + instance + "/databases/" + db})
	if err != nil {
		// Database doesn't exist, or error.
		op, err := admin.CreateDatabase(ctx, &adminpb.CreateDatabaseRequest{
			Parent:          "projects/" + project + "/instances/" + instance,
			CreateStatement: "CREATE DATABASE `" + db + "`",
			ExtraStatements: []string{
				`CREATE TABLE sheep (
							UUID 			STRING(MAX) NOT NULL,
							Count 		INT64
					) PRIMARY KEY (UUID)`,
				`CREATE TABLE sheep_transaction (
							UUID 			STRING(128) NOT NULL,
							Applied 	BOOL
					) PRIMARY KEY (UUID)`,
			},
		})

		if err != nil {
			return err
		}

		_, err = op.Wait(ctx)

		if err != nil {
			return err
		}
	}
	return nil
}

// SetupSpanner initializes the spanner clients.
func SetupSpanner() error {
	ctx := context.Background()
	project := os.Getenv("SHEEP_PROJECT")
	instance := os.Getenv("SHEEP_INSTANCE")
	db := os.Getenv("SHEEP_DATABASE")

	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return err
	}

	// Create the databases if they don't exist.
	err = createSpannerDatabase(ctx, adminClient, project, instance, db)

	if err != nil {
		return err
	}

	dbstr := fmt.Sprintf("projects/%s/instances/%s/databases/%s",
		os.Getenv("SHEEP_PROJECT"),
		instance,
		db)
	client, err := spanner.NewClient(context.Background(), dbstr)

	if err != nil {
		return err
	}

	Spanner = client
	return nil
}
