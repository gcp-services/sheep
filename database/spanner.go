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

func (s *Spanner) Read(msg *Message) error {
	row, err := s.client.Single().ReadRow(
		context.Background(),
		"sheep",
		spanner.Key{msg.Keyspace, msg.Key, msg.Name},
		[]string{"Count"},
	)

	if err != nil {
		return err
	}

	return row.ColumnByName("Count", &msg.Value)
}

func (s *Spanner) Save(message *Message) error {
	ctx := context.WithValue(context.Background(), contextKey("message"), message)
	_, err := s.client.ReadWriteTransaction(ctx, s.doSave)
	return err
}

// Here's where the magic happens. Save out message!
func (s *Spanner) doSave(ctx context.Context, rw *spanner.ReadWriteTransaction) error {
	msg := ctx.Value(contextKey("message")).(*Message)

	// TODO: Since these tables are interleaved, a statement + join
	// will work better here.

	stmt := spanner.NewStatement(`
  SELECT a.Count,
		(SELECT b.UUID
     FROM sheep_transaction AS b
     WHERE b.Keyspace=@Keyspace
     AND b.Key = @Key
     AND b.Name = @Name
     AND b.UUID = @UUID
     ) as UUID
  FROM sheep as a
  WHERE a.Keyspace=@Keyspace
  AND a.Key=@Key
  AND a.Name=@Name
	`)

	stmt.Params["Keyspace"] = msg.Keyspace
	stmt.Params["Key"] = msg.Key
	stmt.Params["Name"] = msg.Name
	stmt.Params["UUID"] = msg.UUID

	iter := rw.Query(ctx, stmt)
	row, err := iter.Next()
	iter.Stop()

	// Let's check and see if our column exists, and if this UUID has been written...
	var uuid spanner.NullString
	var move int64

	if err != nil {
		// If we have a real error, bail.
		if spanner.ErrCode(err) != codes.NotFound {
			return err
		}
		// Not found, which means a new counter we've never seen, so we skip
		// all further checks and exit if here.
	} else {
		// Try to get our UUID
		err = row.ColumnByName("UUID", &uuid)
		// Real error, bail.
		if err != nil {
			return err
		}
		// If the UUID exists in the database, bail, the operation has already been
		// applied.
		if uuid.Valid {
			return nil
		}
		// Get the count.
		err = row.ColumnByName("Count", &move)
		if err != nil {
			return err
		}
	}

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
			Code: codes.InvalidArgument,
			Desc: "Invalid operation sent from message '" + msg.Operation + "', aborting transaction!",
		}
	}

	// Build our mutation...
	m := []*spanner.Mutation{
		spanner.InsertOrUpdate(
			"sheep",
			[]string{"Keyspace", "Key", "Name", "Count"},
			[]interface{}{msg.Keyspace, msg.Key, msg.Name, move},
		),
		spanner.InsertOrUpdate(
			"sheep_transaction",
			[]string{"Keyspace", "Key", "Name", "UUID", "Applied"},
			[]interface{}{msg.Keyspace, msg.Key, msg.Name, msg.UUID, true}),
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
							Keyspace 	STRING(MAX) NOT NULL,
							Key 			STRING(MAX) NOT NULL,
							Name			STRING(MAX) NOT NULL,
							UUID 			STRING(128) NOT NULL,
							Applied 	BOOL
					) PRIMARY KEY (Keyspace, Key, Name, UUID),
						INTERLEAVE IN PARENT sheep ON DELETE CASCADE`,
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
