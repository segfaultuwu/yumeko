package bot

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/segfaultuwu/yumeko/internal/services"
)

func (b *Bot) onUserJoined(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	if m.User == nil {
		return
	}

	settings, err := services.GetGuildSettings(b.DB, m.GuildID)
	if err != nil {
		log.Println("get guild settings:", err)
		return
	}

	if !settings.WelcomeEnabled {
		return
	}

	if settings.WelcomeChannelID == "" {
		log.Println("welcome: welcome_channel_id is empty")
		return
	}

	guildName := "the server"

	guild, err := s.Guild(m.GuildID)
	if err == nil && guild != nil && guild.Name != "" {
		guildName = guild.Name
	}

	content := settings.WelcomeMessage
	content = strings.ReplaceAll(content, "{user_id}", m.User.ID)
	content = strings.ReplaceAll(content, "{username}", m.User.Username)
	content = strings.ReplaceAll(content, "{server_name}", guildName)

	_, err = s.ChannelMessageSend(settings.WelcomeChannelID, content)
	if err != nil {
		log.Println("welcome message:", err)
	}

	if settings.AutoroleID != "" {
		err := s.GuildMemberRoleAdd(m.GuildID, m.User.ID, settings.AutoroleID)
		if err != nil {
			log.Println("autorole:", err)
		}
	}
}
