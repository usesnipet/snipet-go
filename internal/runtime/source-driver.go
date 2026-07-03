package runtime

import (
	"context"
	"time"

	"github.com/usesnipet/snipet/internal/util"
)

type SourceItem struct {
	ID           string
	Name         string       `gorm:"type:text" json:"name"`
	Hash         string       `gorm:"type:varchar(128);index" json:"hash"`
	Metadata     util.JSONMap `gorm:"type:jsonb" json:"metadata"`
	LastModified *time.Time   `json:"last_modified,omitempty"`
}

type SourceDriver interface {
	Scan(ctx context.Context, config util.JSONMap, take *int, skip *int) ([]SourceItem, error)
	TestConnection(ctx context.Context, config util.JSONMap) error
	GetConfigurationSchema(ctx context.Context) (util.JSONMap, error)
}
