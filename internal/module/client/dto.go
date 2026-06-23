package client

type CreateClientDTO struct {
	Name string `json:"name" validate:"required,max=255"`
}

type UpdateClientDTO struct {
	Name *string `json:"name" validate:"omitempty,max=255"`
}
