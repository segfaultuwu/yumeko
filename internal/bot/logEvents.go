package bot

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/segfaultuwu/yumeko/internal/services"
)

func (b *Bot) onMessageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
	if m.GuildID == "" {
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Message deleted",
		Description: "A message was deleted.",
		Color:       0xff4f8b,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Channel",
				Value:  "<#" + m.ChannelID + ">",
				Inline: true,
			},
			{
				Name:   "Message ID",
				Value:  "`" + m.ID + "`",
				Inline: true,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	if err := services.SendGuildLog(b.DB, s, m.GuildID, embed); err != nil {
		log.Println("message delete log:", err)
	}
}

func (b *Bot) onMessageUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) {
	if m.GuildID == "" || m.Author == nil || m.Author.Bot {
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Message edited",
		Description: "A message was edited in <#" + m.ChannelID + ">.",
		Color:       0x9b5cff,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Author",
				Value:  "<@" + m.Author.ID + ">",
				Inline: true,
			},
			{
				Name:   "Message ID",
				Value:  "`" + m.ID + "`",
				Inline: true,
			},
			{
				Name:   "New content",
				Value:  safeEmbedValue(m.Content),
				Inline: false,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	if err := services.SendGuildLog(b.DB, s, m.GuildID, embed); err != nil {
		log.Println("message update log:", err)
	}
}

func (b *Bot) onUserLeft(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
	if m.User == nil {
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Member left",
		Description: "<@" + m.User.ID + "> left the server.",
		Color:       0xff4f8b,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "User",
				Value:  m.User.Username + " (`" + m.User.ID + "`)",
				Inline: false,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	if err := services.SendGuildLog(b.DB, s, m.GuildID, embed); err != nil {
		log.Println("member remove log:", err)
	}
}

func safeEmbedValue(value string) string {
	if value == "" {
		return "`empty`"
	}

	if len(value) > 1000 {
		return value[:1000] + "..."
	}

	return value
}
