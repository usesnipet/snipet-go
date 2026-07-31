package llm

import "github.com/usesnipet/snipet/pkg/msg"

type Prompt struct {
	System   string
	Messages []msg.Message
}
type PromptOption func(*Prompt)

func WithSystem(system string) PromptOption {
	return func(prompt *Prompt) {
		prompt.System = system
	}
}

func WithMessages(messages []msg.Message) PromptOption {
	return func(prompt *Prompt) {
		prompt.Messages = messages
	}
}

func NewPrompt(options ...PromptOption) Prompt {
	prompt := Prompt{}
	for _, option := range options {
		option(&prompt)
	}
	return prompt
}
