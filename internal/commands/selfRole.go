package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/segfaultuwu/yumeko/internal/services"
)

type SelfRoleCommand struct{}

func init() {
	Register(SelfRoleCommand{})
}

func (SelfRoleCommand) Data() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "selfrole",
		Description: "Manage self-assignable roles",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "list",
				Description: "List available self roles",
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "give",
				Description: "Give yourself a self role",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionRole,
						Name:        "role",
						Description: "Role to give yourself",
						Required:    true,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "take",
				Description: "Remove a self role from yourself",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionRole,
						Name:        "role",
						Description: "Role to remove from yourself",
						Required:    true,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "add",
				Description: "Add role to self roles",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionRole,
						Name:        "role",
						Description: "Role to make self-assignable",
						Required:    true,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "remove",
				Description: "Remove role from self roles",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionRole,
						Name:        "role",
						Description: "Role to remove from self roles",
						Required:    true,
					},
				},
			},
		},
	}
}

func (SelfRoleCommand) Execute(ctx Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.GuildID == "" {
		respondEphemeral(s, i, "This command only works on servers.")
		return
	}

	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		respondEphemeral(s, i, "Missing subcommand.")
		return
	}

	sub := options[0]

	switch sub.Name {
	case "list":
		selfRoleList(ctx, s, i)

	case "give":
		role := sub.Options[0].RoleValue(s, i.GuildID)
		if role == nil {
			respondEphemeral(s, i, "Invalid role.")
			return
		}

		selfRoleGive(ctx, s, i, role.ID)

	case "take":
		role := sub.Options[0].RoleValue(s, i.GuildID)
		if role == nil {
			respondEphemeral(s, i, "Invalid role.")
			return
		}

		selfRoleTake(ctx, s, i, role.ID)

	case "add":
		if !hasManageGuild(i.Member) {
			respondEphemeral(s, i, "You need `Manage Server` permission.")
			return
		}

		role := sub.Options[0].RoleValue(s, i.GuildID)
		if role == nil {
			respondEphemeral(s, i, "Invalid role.")
			return
		}

		if err := services.AddSelfRole(ctx.DB, i.GuildID, role.ID); err != nil {
			respondEphemeral(s, i, "Failed to add self role.")
			return
		}

		respondEphemeral(s, i, "Added <@&"+role.ID+"> to self roles.")

	case "remove":
		if !hasManageGuild(i.Member) {
			respondEphemeral(s, i, "You need `Manage Server` permission.")
			return
		}

		role := sub.Options[0].RoleValue(s, i.GuildID)
		if role == nil {
			respondEphemeral(s, i, "Invalid role.")
			return
		}

		if err := services.RemoveSelfRole(ctx.DB, i.GuildID, role.ID); err != nil {
			respondEphemeral(s, i, "Failed to remove self role.")
			return
		}

		respondEphemeral(s, i, "Removed <@&"+role.ID+"> from self roles.")

	default:
		respondEphemeral(s, i, "Unknown subcommand.")
	}
}

func selfRoleList(ctx Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	roles, err := services.ListSelfRoles(ctx.DB, i.GuildID)
	if err != nil {
		respondEphemeral(s, i, "Failed to list self roles.")
		return
	}

	if len(roles) == 0 {
		respondEphemeral(s, i, "No self roles configured.")
		return
	}

	var builder strings.Builder
	builder.WriteString("**Available self roles:**\n\n")

	for _, roleID := range roles {
		builder.WriteString("- <@&")
		builder.WriteString(roleID)
		builder.WriteString(">\n")
	}

	respondEphemeral(s, i, builder.String())
}

func selfRoleGive(ctx Context, s *discordgo.Session, i *discordgo.InteractionCreate, roleID string) {
	user := interactionUser(i)
	if user == nil {
		respondEphemeral(s, i, "Could not detect user.")
		return
	}

	allowed, err := services.IsSelfRole(ctx.DB, i.GuildID, roleID)
	if err != nil {
		respondEphemeral(s, i, "Failed to check self role.")
		return
	}

	if !allowed {
		respondEphemeral(s, i, "This role is not self-assignable.")
		return
	}

	if err := s.GuildMemberRoleAdd(i.GuildID, user.ID, roleID); err != nil {
		respondEphemeral(s, i, "Failed to give role. Check bot permissions and role hierarchy.")
		return
	}

	respondEphemeral(s, i, "Added role <@&"+roleID+"> to you.")
}

func selfRoleTake(ctx Context, s *discordgo.Session, i *discordgo.InteractionCreate, roleID string) {
	user := interactionUser(i)
	if user == nil {
		respondEphemeral(s, i, "Could not detect user.")
		return
	}

	allowed, err := services.IsSelfRole(ctx.DB, i.GuildID, roleID)
	if err != nil {
		respondEphemeral(s, i, "Failed to check self role.")
		return
	}

	if !allowed {
		respondEphemeral(s, i, "This role is not self-assignable.")
		return
	}

	if err := s.GuildMemberRoleRemove(i.GuildID, user.ID, roleID); err != nil {
		respondEphemeral(s, i, "Failed to remove role. Check bot permissions and role hierarchy.")
		return
	}

	respondEphemeral(s, i, "Removed role <@&"+roleID+"> from you.")
}

func hasManageGuild(member *discordgo.Member) bool {
	if member == nil {
		return false
	}

	return member.Permissions&discordgo.PermissionManageGuild == discordgo.PermissionManageGuild
}
