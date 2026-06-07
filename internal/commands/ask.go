package commands

import (
	"context"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/segfaultuwu/yumeko/internal/services"
	"github.com/segfaultuwu/yumeko/internal/tools"
)

type AskCommand struct{}

func init() {
	Register(AskCommand{})
}

func (c AskCommand) Data() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "ask",
		Description: "Ask AI a question",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "prompt",
				Description: "Question for AI",
				Required:    true,
			},
		},
	}
}

func (c AskCommand) Execute(ctx Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options

	if len(options) == 0 {
		respond(s, i, "❌ Missing prompt.")
		return
	}

	prompt := strings.TrimSpace(options[0].StringValue())
	if prompt == "" {
		respond(s, i, "❌ Prompt cannot be empty.")
		return
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		log.Println("defer ask:", err)
		return
	}

	ai := services.NewMistralService(ctx.Config.Ai.MistralAPIKey)

	toolRegistry := tools.NewRegistry()

	answer, err := ai.AskWithTools(
		context.Background(),
		prompt,
		toolRegistry,
	)
	if err != nil {
		log.Println("mistral ask:", err)

		_, _ = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "❌ AI error: " + err.Error(),
		})
		return
	}

	_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: answer,
	})
	if err != nil {
		log.Println("ask followup:", err)
	}
}
