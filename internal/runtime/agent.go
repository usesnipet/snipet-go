package runtime

type LLMConfig Configuration

type Agent struct {
	Name         string
	Description  string
	Instructions string
	LLMs         []LLMConfig
}

func NewAgent(name string, description string, instructions string, llms []LLMConfig) *Agent {
	return &Agent{
		Name:         name,
		Description:  description,
		Instructions: instructions,
		LLMs:         llms,
	}
}
