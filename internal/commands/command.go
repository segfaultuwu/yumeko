package commands

import (
	"database/sql"

	"github.com/bwmarrin/discordgo"
	"github.com/segfaultuwu/yumeko/internal/config"
)

type Context struct {
	DB     *sql.DB
	Config config.Config
}

type Command interface {
	Data() *discordgo.ApplicationCommand
	Execute(ctx Context, s *discordgo.Session, i *discordgo.InteractionCreate)
}
