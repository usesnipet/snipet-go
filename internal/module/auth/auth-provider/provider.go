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
		clientCode string,
		clientConfig *model.ClientConfig,
	) error
	Authenticate(
		ctx context.Context,
		clientCode string,
		clientConfig *model.ClientConfig,
		req *http.Request,
	) (*Identity, error)
}
