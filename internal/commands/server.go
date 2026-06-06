package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type ServerCommand struct{}

func (ServerCommand) Data() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "server",
		Description: "Show server info",
	}
}

func (ServerCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.GuildID == "" {
		respondEphemeral(s, i, "This command only works on servers.")
		return
	}

	guild, err := s.GuildWithCounts(i.GuildID)
	if err != nil {
		respondEphemeral(s, i, "Failed to fetch server info.")
		return
	}

	content := fmt.Sprintf(
		"**%s**\nID: `%s`\nMembers: `%d`\nOwner: <@%s>",
		guild.Name,
		guild.ID,
		guild.ApproximateMemberCount,
		guild.OwnerID,
	)

	respond(s, i, content)
}
