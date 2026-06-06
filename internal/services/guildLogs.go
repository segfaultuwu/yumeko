package services

import (
	"database/sql"

	"github.com/bwmarrin/discordgo"
)

func SendGuildLog(db *sql.DB, s *discordgo.Session, guildID string, embed *discordgo.MessageEmbed) error {
	settings, err := GetGuildSettings(db, guildID)
	if err != nil {
		return err
	}

	if settings.LogChannelID == "" {
		return nil
	}

	_, err = s.ChannelMessageSendEmbed(settings.LogChannelID, embed)
	return err
}
