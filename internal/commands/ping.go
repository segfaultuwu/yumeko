package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type PingCommand struct{}

func init() {
	Register(PingCommand{})
}

func (PingCommand) Data() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "Check bot latency",
	}
}

func (PingCommand) Execute(ctx Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	latency := s.HeartbeatLatency().Milliseconds()
	respond(s, i, fmt.Sprintf("🏓 Pong! `%dms`", latency))
}
