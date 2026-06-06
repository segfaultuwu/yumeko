package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type TodoCommand struct{}

func init() {
	Register(TodoCommand{})
}

func (TodoCommand) Data() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "todo",
		Description: "Manage your todo list",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "add",
				Description: "Add todo",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "content",
						Description: "Todo content",
						Required:    true,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "list",
				Description: "List todos",
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "done",
				Description: "Mark todo as done",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "id",
						Description: "Todo ID",
						Required:    true,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "delete",
				Description: "Delete todo",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "id",
						Description: "Todo ID",
						Required:    true,
					},
				},
			},
		},
	}
}

func (TodoCommand) Execute(ctx Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		respondEphemeral(s, i, "Missing todo subcommand.")
		return
	}

	switch options[0].Name {
	case "add":
		todoAdd(ctx, s, i, options[0])
	case "list":
		todoList(ctx, s, i)
	case "done":
		todoDone(ctx, s, i, options[0])
	case "delete":
		todoDelete(ctx, s, i, options[0])
	default:
		respondEphemeral(s, i, "Unknown todo subcommand.")
	}
}

func todoAdd(ctx Context, s *discordgo.Session, i *discordgo.InteractionCreate, option *discordgo.ApplicationCommandInteractionDataOption) {
	user := interactionUser(i)
	if user == nil {
		respondEphemeral(s, i, "Could not detect user.")
		return
	}

	content := option.Options[0].StringValue()

	_, err := ctx.DB.Exec(
		"INSERT INTO todos (user_id, guild_id, content) VALUES (?, ?, ?)",
		user.ID,
		i.GuildID,
		content,
	)
	if err != nil {
		respondEphemeral(s, i, "Failed to add todo.")
		return
	}

	respondEphemeral(s, i, "Added todo: `"+content+"`")
}

func todoList(ctx Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	user := interactionUser(i)
	if user == nil {
		respondEphemeral(s, i, "Could not detect user.")
		return
	}

	rows, err := ctx.DB.Query(
		"SELECT id, content, done FROM todos WHERE user_id = ? AND guild_id = ? ORDER BY id DESC LIMIT 15",
		user.ID,
		i.GuildID,
	)
	if err != nil {
		respondEphemeral(s, i, "Failed to list todos.")
		return
	}
	defer rows.Close()

	var builder strings.Builder
	builder.WriteString("**Your todos:**\n")

	count := 0

	for rows.Next() {
		var id int64
		var content string
		var done int

		if err := rows.Scan(&id, &content, &done); err != nil {
			respondEphemeral(s, i, "Failed to read todos.")
			return
		}

		status := "⬜"
		if done == 1 {
			status = "✅"
		}

		builder.WriteString(fmt.Sprintf("%s `%d` %s\n", status, id, content))
		count++
	}

	if count == 0 {
		respondEphemeral(s, i, "You have no todos.")
		return
	}

	respondEphemeral(s, i, builder.String())
}

func todoDone(ctx Context, s *discordgo.Session, i *discordgo.InteractionCreate, option *discordgo.ApplicationCommandInteractionDataOption) {
	user := interactionUser(i)
	if user == nil {
		respondEphemeral(s, i, "Could not detect user.")
		return
	}

	id := option.Options[0].IntValue()

	result, err := ctx.DB.Exec(
		"UPDATE todos SET done = 1 WHERE id = ? AND user_id = ? AND guild_id = ?",
		id,
		user.ID,
		i.GuildID,
	)
	if err != nil {
		respondEphemeral(s, i, "Failed to update todo.")
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		respondEphemeral(s, i, "Todo not found.")
		return
	}

	respondEphemeral(s, i, fmt.Sprintf("Marked todo `%d` as done.", id))
}

func todoDelete(ctx Context, s *discordgo.Session, i *discordgo.InteractionCreate, option *discordgo.ApplicationCommandInteractionDataOption) {
	user := interactionUser(i)
	if user == nil {
		respondEphemeral(s, i, "Could not detect user.")
		return
	}

	id := option.Options[0].IntValue()

	result, err := ctx.DB.Exec(
		"DELETE FROM todos WHERE id = ? AND user_id = ? AND guild_id = ?",
		id,
		user.ID,
		i.GuildID,
	)
	if err != nil {
		respondEphemeral(s, i, "Failed to delete todo.")
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		respondEphemeral(s, i, "Todo not found.")
		return
	}

	respondEphemeral(s, i, fmt.Sprintf("Deleted todo `%d`.", id))
}
