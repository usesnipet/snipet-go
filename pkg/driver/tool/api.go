package tool

import (
	"context"

	"github.com/usesnipet/snipet/internal/util"
)

type API struct {
	TestConnection func(ctx context.Context, config util.JSONMap) error

	Call func(ctx context.Context, call Call) (Result, error)
}

func WithAPI(api API) Option {
	return func(o *options) {
		o.api = api
	}
}
