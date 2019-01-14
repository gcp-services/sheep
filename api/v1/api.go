package v1

import (
	"context"
)

type API struct {
}

func (a *API) Get(ctx context.Context, in *Counter) (*Result, error) {
	return &Result{
		Value: 10,
	}, nil
}

func (a *API) Update(ctx context.Context, in *Counter) (*Result, error) {
	return nil, nil
}

func (a *API) Delete(ctx context.Context, in *Counter) (*Result, error) {
	return nil, nil
}
