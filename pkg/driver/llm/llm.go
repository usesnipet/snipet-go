package llm

import (
	"context"

	"github.com/usesnipet/snipet/internal/util"
	"github.com/usesnipet/snipet/pkg/driver"
)

type Driver interface {
	driver.IDriver

	Generate(
		ctx context.Context,
		config util.JSONMap,
		instructions string,
		messages []Message,
	) (Message, error)
}
