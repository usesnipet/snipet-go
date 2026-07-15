package agentloop

import (
	"context"
	"fmt"
)

// ToolSource is a source of tools (function handlers, knowledge/RAG, MCP servers, etc.).
type ToolSource interface {
	Tools() []Tool
	Execute(ctx context.Context, call ToolCall) (*ToolResult, error)
}

// ToolProvider aggregates tool sources and routes execution by tool key.
type ToolProvider interface {
	GetTools() []Tool
	ExecuteTool(ctx context.Context, call ToolCall) (*ToolResult, error)
}

// CompositeToolProvider aggregates multiple ToolSources and routes ExecuteTool by key.
type CompositeToolProvider struct {
	sources []ToolSource
	byKey   map[string]ToolSource
	tools   []Tool
}

func NewCompositeToolProvider(sources ...ToolSource) (*CompositeToolProvider, error) {
	p := &CompositeToolProvider{
		sources: sources,
		byKey:   make(map[string]ToolSource),
	}

	for _, source := range sources {
		for _, tool := range source.Tools() {
			if _, exists := p.byKey[tool.Key]; exists {
				return nil, fmt.Errorf("duplicate tool key %q", tool.Key)
			}
			p.byKey[tool.Key] = source
			p.tools = append(p.tools, tool)
		}
	}

	return p, nil
}

func (p *CompositeToolProvider) GetTools() []Tool {
	return p.tools
}

func (p *CompositeToolProvider) ExecuteTool(ctx context.Context, call ToolCall) (*ToolResult, error) {
	source, ok := p.byKey[call.Key]
	if !ok {
		return nil, fmt.Errorf("unknown tool key %q", call.Key)
	}
	return source.Execute(ctx, call)
}
