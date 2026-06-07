package tools

import (
	"context"
	"time"
)

type TimeTool struct{}

func (TimeTool) Name() string {
	return "time"
}

func (TimeTool) Description() string {
	return "Returns current server time."
}

func (TimeTool) Params() map[string]string {
	return map[string]string{}
}

func (TimeTool) Execute(ctx context.Context, args map[string]any) (string, error) {
	return time.Now().Format(time.RFC3339), nil
}
