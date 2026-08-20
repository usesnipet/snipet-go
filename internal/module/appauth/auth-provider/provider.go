package auth_provider

import (
	"context"
	"net/http"

	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/pkg/jsonx"
)

type ProviderName string

type Identity struct {
	Provider   ProviderName  `json:"provider"`
	ExternalID string        `json:"external_id"`
	Name       *string       `json:"name"`
	Metadata   jsonx.JSONMap `json:"metadata"`
}

type IProvider interface {
	Name() ProviderName
	Validate(
		ctx context.Context,
		appCode string,
		authConfig *model.AppAuthConfig,
	) error
	Authenticate(
		ctx context.Context,
		appCode string,
		authConfig *model.AppAuthConfig,
		req *http.Request,
	) (*Identity, error)
}
