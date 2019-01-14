package v1

import (
	"context"

	"cloud.google.com/go/spanner"
	"github.com/Cidan/sheep/database"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type API struct {
	Stream   database.Stream
	Database database.Database
}

// TODO: Error defs
// TODO: database accepts proto message?
func (a *API) Get(ctx context.Context, in *Counter) (*Result, error) {
	msg := &database.Message{
		Keyspace: in.GetKeyspace(),
		Key:      in.GetKey(),
		Name:     in.GetKey(),
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

func (a *API) Update(ctx context.Context, in *Counter) (*Result, error) {
	return nil, nil
}

func (a *API) Delete(ctx context.Context, in *Counter) (*Result, error) {
	return nil, nil
}
