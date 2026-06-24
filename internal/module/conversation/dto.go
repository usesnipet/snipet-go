package conversation

type CreateConversationDTO struct {
	MemoryID string         `json:"memory_id" validate:"required,uuid"`
	BotID    string         `json:"bot_id" validate:"required,uuid"`
	Metadata map[string]any `json:"metadata" validate:"omitempty"`
}
