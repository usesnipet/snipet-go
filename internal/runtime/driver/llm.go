package driver

import (
	"context"

	"github.com/usesnipet/snipet/internal/runtime/transport"
	"github.com/usesnipet/snipet/internal/util"
)

type ILLM interface {
	IDriver

	Generate(
		ctx context.Context,
		config util.JSONMap,
		instructions string,
		messages []transport.Message,
	) (transport.Message, error)
}
