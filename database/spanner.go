package database

import (
	"context"
	"fmt"

	"cloud.google.com/go/spanner"
	database "cloud.google.com/go/spanner/admin/database/apiv1"
	adminpb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
	"google.golang.org/grpc/codes"
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

func (s *Spanner) Save(message *Message) error {
	ctx := context.WithValue(context.Background(), contextKey("message"), message)
	_, err := s.client.ReadWriteTransaction(ctx, s.doSave)
	return err
}

// Here's where the magic happens. Save out message!
func (s *Spanner) doSave(ctx context.Context, rw *spanner.ReadWriteTransaction) error {
	msg := ctx.Value(contextKey("message")).(*Message)

	// First, let's check and see if our message has been written.
	row, err := rw.ReadRow(context.Background(), "sheep_transaction", spanner.Key{msg.UUID}, []string{"applied"})
	if err != nil {
		if spanner.ErrCode(err) != codes.NotFound {
			return err
		}
	} else {
		var ap bool
		err = row.ColumnByName("Applied", &ap)
		if err != nil {
			return err
		}
		if ap {
			return nil
		}
	}

	// Let's get our current count
	row, err = rw.ReadRow(context.Background(), "sheep", spanner.Key{msg.Keyspace, msg.Key, msg.Name}, []string{"Count"})
	if err != nil {
		return err
	}
	var move int64
	row.ColumnByName("Count", &move)

	// Now we'll do our operation.
	switch msg.Operation {
	case "incr":
		move++
	case "decr":
		move--
	case "set":
		move = msg.Value
	default:
		return &spanner.Error{
			Desc: "Invalid operation sent from message, aborting transaction!",
		}
	}

	// Build our mutation...
	m := []*spanner.Mutation{
		spanner.InsertOrUpdate(
			"sheep_transaction",
			[]string{"UUID", "Applied"},
			[]interface{}{msg.UUID, true}),
		spanner.InsertOrUpdate(
			"sheep",
			[]string{"Keyspace", "Key", "Name", "Count"},
			[]interface{}{msg.Keyspace, msg.Key, msg.Name, move},
		),
	}

	// ...and write!
	return rw.BufferWrite(m)

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
							Keyspace 	STRING(MAX) NOT NULL,
							Key 			STRING(MAX) NOT NULL,
							Name			STRING(MAX) NOT NULL,
							Count 		INT64
					) PRIMARY KEY (Keyspace, Key, Name)`,
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
