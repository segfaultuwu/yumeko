package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type UserCommand struct{}

func init() {
	Register(UserCommand{})
}

func (UserCommand) Data() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "user",
		Description: "Show user info",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "target",
				Description: "User to inspect",
				Required:    false,
			},
		},
	}
}

func (UserCommand) Execute(ctx Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	user := interactionUser(i)

	options := i.ApplicationCommandData().Options
	if len(options) > 0 {
		if options[0].UserValue(s) != nil {
			user = options[0].UserValue(s)
		}
	}

	if user == nil {
		respondEphemeral(s, i, "Could not detect user.")
		return
	}

	content := fmt.Sprintf(
		"**%s**\nID: `%s`\nBot: `%v`\nAvatar: %s",
		user.Username,
		user.ID,
		user.Bot,
		user.AvatarURL("1024"),
	)

	respond(s, i, content)
}
