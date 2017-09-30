package database

import (
	"context"
	"fmt"

	"cloud.google.com/go/spanner"
	database "cloud.google.com/go/spanner/admin/database/apiv1"
	adminpb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
)

type Spanner struct {
	Database
	client *spanner.Client
	admin  *database.DatabaseAdminClient
}

// SetupSpanner initializes the spanner clients.
func NewSpanner(project, instance, db string) (*Spanner, error) {
	ctx := context.Background()
	sp := &Spanner{}

	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return nil, err
	}

	sp.admin = adminClient

	// Create the databases if they don't exist.
	err = sp.createSpannerDatabase(ctx, project, instance, db)

	if err != nil {
		return nil, err
	}

	dbstr := fmt.Sprintf("projects/%s/instances/%s/databases/%s",
		project,
		instance,
		db)

	client, err := spanner.NewClient(context.Background(), dbstr)

	if err != nil {
		return nil, err
	}

	sp.client = client
	return sp, err
}

func (s *Spanner) Read() {
}

func (s *Spanner) Save() {
}

func (s *Spanner) createSpannerDatabase(ctx context.Context, project, instance, db string) error {
	// Create our database if it doesn't exist.
	_, err := s.admin.GetDatabase(ctx, &adminpb.GetDatabaseRequest{
		Name: "projects/" + project + "/instances/" + instance + "/databases/" + db})
	if err != nil {
		// Database doesn't exist, or error.
		op, err := s.admin.CreateDatabase(ctx, &adminpb.CreateDatabaseRequest{
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
