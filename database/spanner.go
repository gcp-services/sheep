package database

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"cloud.google.com/go/spanner"
	database "cloud.google.com/go/spanner/admin/database/apiv1"
	"github.com/Cidan/sheep/stats"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
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

	stmt := spanner.NewStatement(`
			SELECT SUM(a.Count) as Count
			FROM sheep as a
			WHERE a.Keyspace=@Keyspace
			AND a.Key=@Key
			AND a.Name=@Name
		`)
	stmt.Params["Keyspace"] = msg.Keyspace
	stmt.Params["Key"] = msg.Key
	stmt.Params["Name"] = msg.Name

	// TODO: Expose this stale time.
	// TODO: Stale time breaks tests.
	iter := s.client.
		Single().
		//WithTimestampBound(spanner.MaxStaleness(5*time.Second)).
		Query(context.Background(), stmt)
	defer iter.Stop()

	row, err := iter.Next()

	if err != nil {
		return err
	}

	var value spanner.NullInt64
	err = row.ColumnByName("Count", &value)

	if err != nil {
		return err
	}

	if value.Valid {
		msg.Value = value.Int64
		return nil
	}

	return &spanner.Error{
		Code: codes.NotFound,
		Desc: "counter not found",
	}
}

func (s *Spanner) Save(message *Message) error {
	ctx := context.WithValue(context.Background(), contextKey("message"), message)
	if _, err := s.client.ReadWriteTransaction(ctx, s.doSave); err != nil {
		stats.Incr("spanner.save.error", 1)
		return err
	}
	stats.Incr("spanner.save.success", 1)
	return nil
}

// Here's where the magic happens. Save our message!
func (s *Spanner) doSave(ctx context.Context, rw *spanner.ReadWriteTransaction) error {
	msg := ctx.Value(contextKey("message")).(*Message)
	shards := viper.GetInt("spanner.shards")
	shard := rand.Intn(shards)

	// First, let's check and see if our message has been written.
	row, err := rw.ReadRow(context.Background(), "sheep_transaction", spanner.Key{msg.Keyspace, msg.Key, msg.Name, msg.UUID}, []string{"UUID"})
	if err != nil {
		if spanner.ErrCode(err) != codes.NotFound {
			return err
		}
	} else if err == nil {
		// We need to return if err is nil, this means
		// the UUID was found.
		return nil
	}

	// Let's get our current count
	var move int64
	row, err = rw.ReadRow(context.Background(), "sheep", spanner.Key{msg.Keyspace, msg.Key, msg.Name, shard}, []string{"Count"})
	if err != nil {
		if spanner.ErrCode(err) != codes.NotFound {
			return err
		}
	} else {
		row.ColumnByName("Count", &move)
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

	m := []*spanner.Mutation{}

	log.Debug().Int("shard", shard).Msg("shard selected for op")
	if msg.Operation == "set" {
		for i := 0; i < shards; i++ {
			m = append(m, spanner.InsertOrUpdate(
				"sheep",
				[]string{"Keyspace", "Key", "Name", "Shard", "Count"},
				[]interface{}{msg.Keyspace, msg.Key, msg.Name, i, move},
			))
		}
	} else {
		m = append(m, spanner.InsertOrUpdate(
			"sheep",
			[]string{"Keyspace", "Key", "Name", "Shard", "Count"},
			[]interface{}{msg.Keyspace, msg.Key, msg.Name, shard, move},
		))
	}

	m = append(m, spanner.InsertOrUpdate(
		"sheep_transaction",
		[]string{"Keyspace", "Key", "Name", "UUID", "Time"},
		[]interface{}{msg.Keyspace, msg.Key, msg.Name, msg.UUID, time.Now()}))

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
							Shard     INT64       NOT NULL,
							Count 		INT64       NOT NULL
					) PRIMARY KEY (Keyspace, Key, Name, Shard)`,
				`CREATE TABLE sheep_transaction (
							Keyspace 	STRING(MAX) NOT NULL,
							Key 			STRING(MAX) NOT NULL,
							Name			STRING(MAX) NOT NULL,
							UUID 			STRING(128) NOT NULL,
							Time      TIMESTAMP   NOT NULL
					) PRIMARY KEY (Keyspace, Key, Name, UUID)`,
				`CREATE TABLE sheep_stats (
					   UUID      STRING(MAX) NOT NULL,
						 Key       STRING(MAX) NOT NULL,
						 Value     FLOAT64     NOT NULL,
						 Hostname  STRING(MAX) NOT NULL,
						 Last      TIMESTAMP   NOT NULL
				 ) PRIMARY KEY (UUID, Key)`,
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
