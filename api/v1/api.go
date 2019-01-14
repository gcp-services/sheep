package v1

import (
	"context"

	"cloud.google.com/go/spanner"
	"github.com/Cidan/sheep/database"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// API struct for v1 of the go API
type API struct {
	Stream   database.Stream
	Database database.Database
}

// Get a counter result
// TODO: Error defs
// TODO: database accepts proto message?
func (a *API) Get(ctx context.Context, in *Counter) (*Result, error) {
	msg := &database.Message{
		Keyspace: in.GetKeyspace(),
		Key:      in.GetKey(),
		Name:     in.GetName(),
	}
	err := a.Database.Read(msg)
	if err != nil {
		if spanner.ErrCode(err) == codes.NotFound {
			return &Result{}, status.Error(codes.NotFound, "Counter not found in database")
		}
		return &Result{}, status.Error(codes.Internal, err.Error())
	}

	return &Result{
		Value: msg.Value,
	}, nil
}

// Update a counter in the database
func (a *API) Update(ctx context.Context, in *Counter) (*Result, error) {
	var err error

	msg := &database.Message{
		Keyspace:  in.GetKeyspace(),
		Key:       in.GetKey(),
		Name:      in.GetName(),
		UUID:      in.GetUuid(),
		Operation: in.GetOperation().String(),
	}

	if in.GetDirect() {
		err = a.Database.Save(msg)
	} else {
		err = a.Stream.Save(msg)
	}
	if err != nil {
		return &Result{}, status.Error(codes.Internal, err.Error())
	}
	return &Result{Value: msg.Value}, nil

}

// Delete a counter
func (a *API) Delete(ctx context.Context, in *Counter) (*Result, error) {
	return nil, nil
}
