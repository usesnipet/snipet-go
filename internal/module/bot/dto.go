package bot

import "github.com/usesnipet/snipet/internal/model"

type CreateBotDTO struct {
	Name          string                 `json:"name" validate:"required,max=255"`
	Description   string                 `json:"description" validate:"max=1000"`
	Configuration model.BotConfiguration `json:"configuration" validate:"required"`
}

type UpdateBotDTO struct {
	Name          *string                 `json:"name" validate:"omitempty,max=255"`
	Description   *string                 `json:"description" validate:"omitempty,max=1000"`
	Configuration *model.BotConfiguration `json:"configuration"`
}

type LinkClientToBotDTO struct {
	ClientCode string `json:"client_code" validate:"required,max=10"`
	BotID      string `json:"bot_id" validate:"required,uuid"`
}
