package tools

import "context"

type PingTool struct{}

func (PingTool) Name() string {
	return "ping"
}

func (PingTool) Description() string {
	return "Checks if tools system works. Returns pong."
}

func (PingTool) Params() map[string]string {
	return map[string]string{}
}

func (PingTool) Execute(ctx context.Context, args map[string]any) (string, error) {
	return "pong", nil
}
