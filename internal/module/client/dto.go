package client

type CreateClientDTO struct {
	Name       string `json:"name" validate:"required,max=255"`
	WebhookURL string `json:"webhook_url" validate:"omitempty,url"`
}

type UpdateClientDTO struct {
	Name       *string `json:"name" validate:"omitempty,max=255"`
	WebhookURL *string `json:"webhook_url" validate:"omitempty,url"`
}
