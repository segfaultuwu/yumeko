package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/segfaultuwu/yumeko/internal/commands"
)

func (b *Bot) registerEvents() {
	b.Session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Ready:", s.State.User.Username)
	})

	b.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}

		name := i.ApplicationCommandData().Name

		handler, ok := b.commandHandlers()[name]
		if !ok {
			RespondEphemeral(s, i, "Unknown command.")
			return
		}

		ctx := commands.Context{
			DB:     b.DB,
			Config: b.Config,
		}

		handler.Execute(ctx, s, i)
	})

	b.Session.AddHandler(func(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
		b.onUserJoined(s, m)
	})

	b.Session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageDelete) {
		b.onMessageDelete(s, m)
	})

	b.Session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageUpdate) {
		b.onMessageUpdate(s, m)
	})

	b.Session.AddHandler(func(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
		b.onUserLeft(s, m)
	})
}
