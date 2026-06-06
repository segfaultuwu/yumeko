package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type PollCommand struct{}

func init() {
	Register(PollCommand{})
}

func (PollCommand) Data() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "poll",
		Description: "Create yes/no poll",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "question",
				Description: "Poll question",
				Required:    true,
			},
		},
	}
}

func (PollCommand) Execute(ctx Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		respondEphemeral(s, i, "Missing poll question.")
		return
	}

	question := options[0].StringValue()

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("📊 **Poll:** %s\n\n✅ = Yes\n❌ = No", question),
		},
	})
	if err != nil {
		return
	}

	msg, err := s.InteractionResponse(i.Interaction)
	if err != nil {
		return
	}

	_ = s.MessageReactionAdd(i.ChannelID, msg.ID, "✅")
	_ = s.MessageReactionAdd(i.ChannelID, msg.ID, "❌")
}
