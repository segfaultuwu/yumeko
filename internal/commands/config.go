package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/segfaultuwu/yumeko/internal/services"
)

type ConfigCommand struct{}

func init() {
	Register(ConfigCommand{})
}

func (ConfigCommand) Data() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "config",
		Description: "Manage server config",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "show",
				Description: "Show current server config",
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "welcome-channel",
				Description: "Set welcome channel",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:         discordgo.ApplicationCommandOptionChannel,
						Name:         "channel",
						Description:  "Welcome channel",
						Required:     true,
						ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "welcome-enabled",
				Description: "Enable or disable welcome messages",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionBoolean,
						Name:        "enabled",
						Description: "Enabled",
						Required:    true,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "welcome-message",
				Description: "Set welcome message",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "message",
						Description: "Message with {user_id}, {username}, {server_name}",
						Required:    true,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "log-channel",
				Description: "Set moderation log channel",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:         discordgo.ApplicationCommandOptionChannel,
						Name:         "channel",
						Description:  "Log channel",
						Required:     true,
						ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "autorole",
				Description: "Set autorole",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionRole,
						Name:        "role",
						Description: "Role to give on join",
						Required:    true,
					},
				},
			},
		},
	}
}

func (ConfigCommand) Execute(ctx Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := deferEphemeral(s, i); err != nil {
		return
	}

	if i.GuildID == "" {
		editResponse(s, i, "This command only works on servers.")
		return
	}

	if i.Member == nil {
		editResponse(s, i, "Could not detect member.")
		return
	}

	if !memberHasPermission(i.Member, discordgo.PermissionManageGuild) {
		editResponse(s, i, "You need `Manage Server` permission.")
		return
	}

	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		editResponse(s, i, "Missing config subcommand.")
		return
	}

	sub := options[0]

	switch sub.Name {
	case "show":
		configShow(ctx, s, i)

	case "welcome-channel":
		if len(sub.Options) == 0 {
			editResponse(s, i, "Missing channel.")
			return
		}

		channel := sub.Options[0].ChannelValue(s)
		if channel == nil {
			editResponse(s, i, "Invalid channel.")
			return
		}

		if err := services.SetWelcomeChannel(ctx.DB, i.GuildID, channel.ID); err != nil {
			editResponse(s, i, "Failed to set welcome channel: `"+err.Error()+"`")
			return
		}

		editResponse(s, i, "Welcome channel set to <#"+channel.ID+">.")

	case "welcome-enabled":
		if len(sub.Options) == 0 {
			editResponse(s, i, "Missing enabled value.")
			return
		}

		enabled := sub.Options[0].BoolValue()

		if err := services.SetWelcomeEnabled(ctx.DB, i.GuildID, enabled); err != nil {
			editResponse(s, i, "Failed to update welcome setting: `"+err.Error()+"`")
			return
		}

		editResponse(s, i, fmt.Sprintf("Welcome enabled: `%v`.", enabled))

	case "welcome-message":
		if len(sub.Options) == 0 {
			editResponse(s, i, "Missing message.")
			return
		}

		message := sub.Options[0].StringValue()

		if err := services.SetWelcomeMessage(ctx.DB, i.GuildID, message); err != nil {
			editResponse(s, i, "Failed to set welcome message: `"+err.Error()+"`")
			return
		}

		editResponse(s, i, "Welcome message updated.")

	case "log-channel":
		if len(sub.Options) == 0 {
			editResponse(s, i, "Missing channel.")
			return
		}

		channel := sub.Options[0].ChannelValue(s)
		if channel == nil {
			editResponse(s, i, "Invalid channel.")
			return
		}

		if err := services.SetLogChannel(ctx.DB, i.GuildID, channel.ID); err != nil {
			editResponse(s, i, "Failed to set log channel: `"+err.Error()+"`")
			return
		}

		editResponse(s, i, "Log channel set to <#"+channel.ID+">.")

	case "autorole":
		if len(sub.Options) == 0 {
			editResponse(s, i, "Missing role.")
			return
		}

		role := sub.Options[0].RoleValue(s, i.GuildID)
		if role == nil {
			editResponse(s, i, "Invalid role.")
			return
		}

		if err := services.SetAutorole(ctx.DB, i.GuildID, role.ID); err != nil {
			editResponse(s, i, "Failed to set autorole: `"+err.Error()+"`")
			return
		}

		editResponse(s, i, "Autorole set to <@&"+role.ID+">.")

	default:
		editResponse(s, i, "Unknown config subcommand.")
	}
}

func configShow(ctx Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	settings, err := services.GetGuildSettings(ctx.DB, i.GuildID)
	if err != nil {
		respondEphemeral(s, i, "Failed to read server config.")
		return
	}

	welcomeChannel := "`not set`"
	if settings.WelcomeChannelID != "" {
		welcomeChannel = "<#" + settings.WelcomeChannelID + ">"
	}

	logChannel := "`not set`"
	if settings.LogChannelID != "" {
		logChannel = "<#" + settings.LogChannelID + ">"
	}

	autorole := "`not set`"
	if settings.AutoroleID != "" {
		autorole = "<@&" + settings.AutoroleID + ">"
	}

	content := fmt.Sprintf(
		"**Server config**\n\nWelcome enabled: `%v`\nWelcome channel: %s\nLog channel: %s\nAutorole: %s\nWelcome message:\n```%s```",
		settings.WelcomeEnabled,
		welcomeChannel,
		logChannel,
		autorole,
		settings.WelcomeMessage,
	)

	respondEphemeral(s, i, content)
}

func memberHasPermission(member *discordgo.Member, permission int64) bool {
	return member.Permissions&permission == permission
}
