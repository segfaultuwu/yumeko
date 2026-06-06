package bot

import (
	"fmt"

	"github.com/segfaultuwu/yumeko/internal/commands"
)

func (b *Bot) commandHandlers() map[string]commands.Command {
	handlers := make(map[string]commands.Command)

	for _, cmd := range commands.All() {
		handlers[cmd.Data().Name] = cmd
	}

	return handlers
}

func (b *Bot) RegisterCommands() error {
	appID := b.Session.State.User.ID

	for _, cmd := range commands.All() {
		data := cmd.Data()

		_, err := b.Session.ApplicationCommandCreate(
			appID,
			b.Config.Bot.GuildID,
			data,
		)
		if err != nil {
			return fmt.Errorf("register command %s: %w", data.Name, err)
		}
	}

	return nil
}
