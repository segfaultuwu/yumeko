package bot

import (
	"context"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/segfaultuwu/yumeko/internal/services"
	"github.com/segfaultuwu/yumeko/internal/tools"
)

func (b *Bot) onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author == nil || m.Author.Bot {
		return
	}

	if s.State == nil || s.State.User == nil {
		return
	}

	botID := s.State.User.ID
	content := strings.TrimSpace(m.Content)

	mentioned := isBotMentioned(m.Message, botID)
	repliedToBot := b.isReplyToBot(s, m.Message, botID)

	if !mentioned && !repliedToBot {
		return
	}

	clean := cleanBotMention(content, botID)

	if clean == "" {
		return
	}

	if strings.EqualFold(clean, "ping") {
		_, err := s.ChannelMessageSendReply(
			m.ChannelID,
			"pong 🏓",
			m.Reference(),
		)
		if err != nil {
			log.Println("send ping reply:", err)
		}

		return
	}

	_ = s.ChannelTyping(m.ChannelID)

	ai := services.NewMistralService(b.Config.Ai.MistralAPIKey)

	toolRegistry := tools.NewRegistry()

	answer, err := ai.AskWithTools(
		context.Background(),
		clean,
		toolRegistry,
	)
	if err != nil {
		log.Println("mistral ask:", err)

		_, _ = s.ChannelMessageSendReply(
			m.ChannelID,
			"❌ AI error: "+err.Error(),
			m.Reference(),
		)

		return
	}

	if len(answer) > 1900 {
		answer = answer[:1900] + "..."
	}

	_, err = s.ChannelMessageSendReply(
		m.ChannelID,
		answer,
		m.Reference(),
	)
	if err != nil {
		log.Println("send ai reply:", err)
	}
}

func isBotMentioned(m *discordgo.Message, botID string) bool {
	for _, user := range m.Mentions {
		if user.ID == botID {
			return true
		}
	}

	return false
}

func cleanBotMention(content string, botID string) string {
	content = strings.ReplaceAll(content, "<@"+botID+">", "")
	content = strings.ReplaceAll(content, "<@!"+botID+">", "")

	return strings.TrimSpace(content)
}

func (b *Bot) isReplyToBot(s *discordgo.Session, m *discordgo.Message, botID string) bool {
	if m.MessageReference == nil || m.MessageReference.MessageID == "" {
		return false
	}

	ref, err := s.ChannelMessage(m.ChannelID, m.MessageReference.MessageID)
	if err != nil {
		log.Println("fetch reply message:", err)
		return false
	}

	return ref.Author != nil && ref.Author.ID == botID
}
