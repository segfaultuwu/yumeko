package services

import (
	"database/sql"
)

type GuildSettings struct {
	GuildID          string
	LogChannelID     string
	WelcomeChannelID string
	AutoroleID       string
	WelcomeEnabled   bool
	WelcomeMessage   string
}

func GetGuildSettings(db *sql.DB, guildID string) (GuildSettings, error) {
	var s GuildSettings
	var welcomeEnabled int

	err := db.QueryRow(`
		SELECT
			guild_id,
			COALESCE(log_channel_id, ''),
			COALESCE(welcome_channel_id, ''),
			COALESCE(autorole_id, ''),
			welcome_enabled,
			welcome_message
		FROM guild_settings
		WHERE guild_id = ?
	`, guildID).Scan(
		&s.GuildID,
		&s.LogChannelID,
		&s.WelcomeChannelID,
		&s.AutoroleID,
		&welcomeEnabled,
		&s.WelcomeMessage,
	)

	if err == sql.ErrNoRows {
		return DefaultGuildSettings(guildID), nil
	}

	if err != nil {
		return s, err
	}

	s.WelcomeEnabled = welcomeEnabled == 1
	return s, nil
}

func DefaultGuildSettings(guildID string) GuildSettings {
	return GuildSettings{
		GuildID:          guildID,
		WelcomeEnabled:   false,
		WelcomeMessage:   "👋 Welcome <@{user_id}> to **{server_name}**!",
		LogChannelID:     "",
		WelcomeChannelID: "",
		AutoroleID:       "",
	}
}

func SetWelcomeChannel(db *sql.DB, guildID, channelID string) error {
	_, err := db.Exec(`
		INSERT INTO guild_settings (guild_id, welcome_channel_id, updated_at)
		VALUES (?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(guild_id) DO UPDATE SET
			welcome_channel_id = excluded.welcome_channel_id,
			updated_at = CURRENT_TIMESTAMP
	`, guildID, channelID)

	return err
}

func SetLogChannel(db *sql.DB, guildID, channelID string) error {
	_, err := db.Exec(`
		INSERT INTO guild_settings (guild_id, log_channel_id, updated_at)
		VALUES (?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(guild_id) DO UPDATE SET
			log_channel_id = excluded.log_channel_id,
			updated_at = CURRENT_TIMESTAMP
	`, guildID, channelID)

	return err
}

func SetAutorole(db *sql.DB, guildID, roleID string) error {
	_, err := db.Exec(`
		INSERT INTO guild_settings (guild_id, autorole_id, updated_at)
		VALUES (?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(guild_id) DO UPDATE SET
			autorole_id = excluded.autorole_id,
			updated_at = CURRENT_TIMESTAMP
	`, guildID, roleID)

	return err
}

func SetWelcomeEnabled(db *sql.DB, guildID string, enabled bool) error {
	value := 0
	if enabled {
		value = 1
	}

	_, err := db.Exec(`
		INSERT INTO guild_settings (guild_id, welcome_enabled, updated_at)
		VALUES (?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(guild_id) DO UPDATE SET
			welcome_enabled = excluded.welcome_enabled,
			updated_at = CURRENT_TIMESTAMP
	`, guildID, value)

	return err
}

func SetWelcomeMessage(db *sql.DB, guildID, message string) error {
	_, err := db.Exec(`
		INSERT INTO guild_settings (guild_id, welcome_message, updated_at)
		VALUES (?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(guild_id) DO UPDATE SET
			welcome_message = excluded.welcome_message,
			updated_at = CURRENT_TIMESTAMP
	`, guildID, message)

	return err
}
