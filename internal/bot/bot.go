package bot

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/segfaultuwu/yumeko/internal/config"
)

type Bot struct {
	Session *discordgo.Session
	Config  config.Config
	DB      *sql.DB
}

func New(cfg config.Config, db *sql.DB) (*Bot, error) {
	session, err := discordgo.New("Bot " + cfg.Bot.Token)
	if err != nil {
		return nil, err
	}

	session.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMembers

	b := &Bot{
		Session: session,
		Config:  cfg,
		DB:      db,
	}

	b.registerEvents()

	return b, nil
}

func (b *Bot) Start() error {
	if err := b.Session.Open(); err != nil {
		return err
	}

	log.Println("Connected as", b.Session.State.User.Username)

	if err := b.RegisterCommands(); err != nil {
		return fmt.Errorf("register commands: %w", err)
	}

	return nil
}

func (b *Bot) Stop() error {
	return b.Session.Close()
}
