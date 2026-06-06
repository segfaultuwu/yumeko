package config

import (
	"errors"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Bot struct {
		Token   string `toml:"token"`
		GuildID string `toml:"guild_id"`
	} `toml:"bot"`
	Guild struct {
		Name string `toml:"name"`
	} `toml:"guild"`
	Database struct {
		Path string `toml:"path"`
	} `toml:"database"`

	Settings struct {
		OwnerID string `toml:"owner_id"`
	} `toml:"settings"`

	Welcome struct {
		Enabled   bool   `toml:"enabled"`
		ChannelID string `toml:"channel_id"`
		Message   string `toml:"message"`
	} `toml:"welcome"`
}

func Load(path string) (Config, error) {
	var cfg Config

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}

	if err := toml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}

	if cfg.Bot.Token == "" {
		return cfg, errors.New("bot token is empty")
	}

	if cfg.Database.Path == "" {
		cfg.Database.Path = "./data/yumeko.db"
	}

	if cfg.Welcome.Message == "" {
		cfg.Welcome.Message = "👋 Welcome <@{user_id}> to **{server_name}**!"
	}

	return cfg, nil
}
