package apikey

import (
	"time"

	"github.com/usesnipet/snipet/internal/model"
)

type CreateAPIKeyDTO struct {
	Name      string     `json:"name" validate:"required,max=255"`
	ExpiresAt *time.Time `json:"expires_at"`
}

type UpdateExpirationDTO struct {
	ExpiresAt *time.Time `json:"expires_at"`
}

type APIKeyWithSecret struct {
	*model.APIKey
	Key string `json:"key"`
}
