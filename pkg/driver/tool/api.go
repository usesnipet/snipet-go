package tool

import (
	"context"

	"github.com/usesnipet/snipet/pkg/jsonx"
)

// API holds the provider methods used by a driver for connectivity checks
// and invoking a tool call.
type API struct {
	TestConnection func(ctx context.Context, config jsonx.JSONMap) error

	Call func(ctx context.Context, call Call) (Result, error)
}

// WithAPI sets the driver's underlying implementation.
func WithAPI(api API) Option {
	return func(o *options) {
		o.api = api
	}
}
