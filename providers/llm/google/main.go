package google

import (
	"context"

	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/provider/llm"
	"github.com/usesnipet/snipet/internal/util"
	"google.golang.org/genai"
)

type GoogleLLMProvider struct {
}

func (p *GoogleLLMProvider) Run(
	ctx context.Context,
	configuration llm.Configuration[Config],
	messages []model.SessionMessage,
) ([]model.SessionMessagePart, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  configuration.Config.APIKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return []model.SessionMessagePart{}, err
	}

	contents := make([]*genai.Content, len(messages))
	for i, message := range messages {
		contents[i] = &genai.Content{
			Role: message.Role,
			Parts: util.Map(
				message.Parts,
				func(part model.SessionMessagePart) *genai.Part {
					return &genai.Part{
						Text: part.Content.(string),
					}
				},
			),
		}
	}
	res, err := client.Models.GenerateContent(ctx, configuration.Model, contents, nil)
	if err != nil {
		return []model.SessionMessagePart{}, err
	}
	return util.Map(
		res.Candidates,
		func(candidate *genai.Candidate) model.SessionMessagePart {
			return model.SessionMessagePart{
				Type:    model.SessionMessagePartTypeText,
				Content: candidate.Content.Parts[0].Text,
			}
		},
	), nil
}

func (p *GoogleLLMProvider) Name() string {
	return "google"
}

func New() llm.Provider[Config] {
	return &GoogleLLMProvider{}
}
