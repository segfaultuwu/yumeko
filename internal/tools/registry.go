package tools

import (
	"context"
	"fmt"
)

type Registry struct {
	tools map[string]Tool
}

func NewRegistry() *Registry {
	r := &Registry{
		tools: make(map[string]Tool),
	}

	r.Register(PingTool{})
	r.Register(TimeTool{})

	return r
}

func (r *Registry) Register(tool Tool) {
	r.tools[tool.Name()] = tool
}

func (r *Registry) Get(name string) (Tool, bool) {
	tool, ok := r.tools[name]
	return tool, ok
}

func (r *Registry) Execute(ctx context.Context, name string, args map[string]any) (string, error) {
	tool, ok := r.Get(name)
	if !ok {
		return "", fmt.Errorf("unknown tool: %s", name)
	}

	return tool.Execute(ctx, args)
}

func (r *Registry) Prompt() string {
	out := "Available tools:\n"

	for _, tool := range r.tools {
		out += fmt.Sprintf(
			"- %s: %s\n",
			tool.Name(),
			tool.Description(),
		)
	}

	return out
}
