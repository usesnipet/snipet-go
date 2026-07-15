package agentloop

import "context"

type LLMProvider interface {
	GenerateResponse(ctx context.Context, agent *Agent, execution *Execution) (*Message, error)
}
