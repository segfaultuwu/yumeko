package commands

import "github.com/bwmarrin/discordgo"

type HelpCommand struct{}

func init() {
	Register(HelpCommand{})
}

func (HelpCommand) Data() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "help",
		Description: "Show Yumeko commands",
	}
}

func (HelpCommand) Execute(ctx Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	content := "**Yumeko commands**\n\n"

	for _, cmd := range All() {
		data := cmd.Data()
		content += "`/" + data.Name + "` - " + data.Description + "\n"
	}

	respondEphemeral(s, i, content)
}
