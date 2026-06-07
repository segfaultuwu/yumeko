package tools

import "context"

type Tool interface {
	Name() string
	Description() string
	Params() map[string]string
	Execute(ctx context.Context, args map[string]any) (string, error)
}
