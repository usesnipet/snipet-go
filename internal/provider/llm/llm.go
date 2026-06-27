package llm

import (
	"context"

	"github.com/usesnipet/snipet/internal/model"
)

type Provider[TConfig any] interface {
	Name() string
	Run(
		ctx context.Context,
		configuration Configuration[TConfig],
		messages []model.SessionMessage,
	) ([]model.SessionMessagePart, error)
}
