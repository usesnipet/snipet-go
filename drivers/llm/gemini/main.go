package gemini

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"time"

	"github.com/usesnipet/snipet/internal/runtime/driver"
	"github.com/usesnipet/snipet/internal/runtime/tool"
	"github.com/usesnipet/snipet/internal/runtime/transport"
	"github.com/usesnipet/snipet/internal/util"
	jsonschema "github.com/usesnipet/snipet/internal/util/json_schema"
	"google.golang.org/genai"
)

//go:embed schema.json
var schemaJSON []byte

type Gemini struct{}

func New() driver.ILLM {
	return &Gemini{}
}

func (g *Gemini) Info() driver.Info {
	schema, _ := jsonschema.Load(schemaJSON)

	return driver.Info{
		Name:                "Gemini",
		Description:         "Gemini is a language model that can generate text, images, and audio.",
		Icon:                "https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png",
		Tags:                []string{"language", "model", "llm"},
		ConfigurationSchema: schema,
	}
}

func (g *Gemini) TestConnection(ctx context.Context, config util.JSONMap) error {
	cfg, err := util.ParseJSONMap[Config](config)
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  cfg.APIKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	_, err = client.Models.List(ctx, &genai.ListModelsConfig{
		PageSize: 1,
	})
	if err != nil {
		return fmt.Errorf("failed to list models: %w", err)
	}

	return nil
}

func (g *Gemini) Generate(
	ctx context.Context,
	config util.JSONMap,
	instructions string,
	messages []transport.Message,
) (message transport.Message, err error) {
	cfg, err := util.ParseJSONMap[Config](config)
	if err != nil {
		return message, fmt.Errorf("failed to parse config: %w", err)
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  cfg.APIKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return message, fmt.Errorf("failed to create client: %w", err)
	}

	contents, err := buildContents(messages)
	if err != nil {
		return message, fmt.Errorf("gemini: build contents: %w", err)
	}

	genCfg := &genai.GenerateContentConfig{
		ResponseMIMEType:   "application/json",
		ResponseJsonSchema: messageSchema,
	}
	if instructions != "" {
		// Empty text is omitted by json omitempty, which leaves an invalid Part with no data.
		genCfg.SystemInstruction = genai.NewContentFromText(instructions, genai.RoleUser)
	}
	if cfg.Temperature != 0 {
		genCfg.Temperature = genai.Ptr(float32(cfg.Temperature))
	}
	if cfg.TopP != 0 {
		genCfg.TopP = genai.Ptr(float32(cfg.TopP))
	}
	if cfg.MaxTokens != 0 {
		genCfg.MaxOutputTokens = int32(cfg.MaxTokens)
	}
	result, err := client.Models.GenerateContent(ctx, cfg.Model, contents, genCfg)
	if err != nil {
		return message, fmt.Errorf("gemini: generate response: %w", err)
	}

	text := result.Text()
	if text == "" {
		return message, fmt.Errorf("gemini: empty response")
	}

	var out Message
	if err := json.Unmarshal([]byte(text), &out); err != nil {
		return message, fmt.Errorf("gemini: parse structured output: %w", err)
	}

	return toAgentMessage(out), nil
}

func buildContents(messages []transport.Message) ([]*genai.Content, error) {
	out := make([]*genai.Content, 0, len(messages))
	for _, message := range messages {
		item, ok, err := toContent(message)
		if err != nil {
			return nil, err
		}
		if ok {
			out = append(out, item)
		}
	}
	return out, nil
}

func toContent(message transport.Message) (*genai.Content, bool, error) {
	switch message.Role {
	case transport.MessageRoleUser:
		return genai.NewContentFromText(message.Content, genai.RoleUser), true, nil
	case transport.MessageRoleSystem:
		return genai.NewContentFromText(message.Content, genai.RoleUser), true, nil
	case transport.MessageRoleAssistant, transport.MessageRoleFinal:
		content := message.Content
		if len(message.ToolCalls) > 0 {
			raw, err := json.Marshal(Message{
				Content:   message.Content,
				ToolCalls: util.Map(message.ToolCalls, fromAgentToolCall),
			})
			if err != nil {
				return nil, false, fmt.Errorf("marshal assistant message: %w", err)
			}
			content = string(raw)
		}
		return genai.NewContentFromText(content, genai.RoleModel), true, nil
	case transport.MessageRoleTool:
		content := message.Content
		if message.ToolResult != nil {
			content = fmt.Sprintf(
				"Tool %q result: %s",
				message.ToolResult.Key,
				message.Content,
			)
		}
		return genai.NewContentFromText(content, genai.RoleUser), true, nil
	default:
		return nil, false, nil
	}
}

func toAgentMessage(out Message) transport.Message {
	toolCalls := make([]tool.Call, 0, len(out.ToolCalls))
	for _, call := range out.ToolCalls {
		toolCalls = append(toolCalls, toAgentToolCall(call))
	}

	role := transport.MessageRoleFinal
	if len(toolCalls) > 0 {
		role = transport.MessageRoleAssistant
	}

	return transport.Message{
		Role:      role,
		Content:   out.Content,
		ToolCalls: toolCalls,
		Timestamp: time.Now(),
	}
}

func toAgentToolCall(call ToolCall) tool.Call {
	var input any
	if call.Input != "" {
		if err := json.Unmarshal([]byte(call.Input), &input); err != nil {
			input = call.Input
		}
	}

	return tool.Call{
		Key:   call.Key,
		Input: input,
	}
}

func fromAgentToolCall(call tool.Call) ToolCall {
	input := ""
	if call.Input != nil {
		if raw, err := json.Marshal(call.Input); err == nil {
			input = string(raw)
		} else {
			input = fmt.Sprint(call.Input)
		}
	}
	return ToolCall{
		Key:   call.Key,
		Input: input,
	}
}
