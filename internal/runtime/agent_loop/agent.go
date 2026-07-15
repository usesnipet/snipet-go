package agentloop

type Agent struct {
	ID           string
	Name         string
	Description  string
	Instructions string
	Tools        []Tool
	LLMs         []any
}

func NewAgent(id string, name string, description string, instructions string, tools []Tool) *Agent {
	return &Agent{
		ID:           id,
		Name:         name,
		Description:  description,
		Instructions: instructions,
		Tools:        tools,
	}
}
