package commands

import "sync"

var (
	mu       sync.RWMutex
	registry []Command
)

func Register(cmd Command) {
	mu.Lock()
	defer mu.Unlock()

	registry = append(registry, cmd)
}

func All() []Command {
	mu.RLock()
	defer mu.RUnlock()

	out := make([]Command, len(registry))
	copy(out, registry)

	return out
}
