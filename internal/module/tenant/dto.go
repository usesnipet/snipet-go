package tenant

import (
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
)

// TenantResponse and TenantsPage exist so swagger annotations in this
// package can reference them without importing internal/model or
// internal/page directly.
type TenantResponse = model.Tenant

type TenantsPage = page.Paginated[model.Tenant]

type CreateTenantDTO struct {
	Name string  `json:"name" validate:"required,max=255"`
	Slug string  `json:"slug" validate:"required,max=255,alphanum_dash"`
	Icon *string `json:"icon" validate:"omitempty,max=255"`
}

type UpdateTenantDTO struct {
	Name *string `json:"name" validate:"omitempty,max=255"`
	Slug *string `json:"slug" validate:"omitempty,max=255,alphanum_dash"`
	Icon *string `json:"icon" validate:"omitempty,max=255"`
}
