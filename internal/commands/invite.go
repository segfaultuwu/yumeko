package commands

import (
	"fmt"
	"net/url"

	"github.com/bwmarrin/discordgo"
)

type InviteCommand struct{}

func init() {
	Register(InviteCommand{})
}

func (InviteCommand) Data() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "invite",
		Description: "Get Yumeko invite link",
	}
}

func (InviteCommand) Execute(ctx Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	clientID := s.State.User.ID
	if clientID == "" {
		respondEphemeral(s, i, "Could not detect bot client ID.")
		return
	}

	// Required permissions:
	// Send Messages
	// Embed Links
	// Add Reactions
	// Use Slash Commands
	// Manage Roles - for autoroles
	permissions := discordgo.PermissionSendMessages |
		discordgo.PermissionEmbedLinks |
		discordgo.PermissionAddReactions |
		discordgo.PermissionManageRoles

	inviteURL := fmt.Sprintf(
		"https://discord.com/oauth2/authorize?client_id=%s&permissions=%d&scope=%s",
		url.QueryEscape(clientID),
		permissions,
		url.QueryEscape("bot applications.commands"),
	)

	content := fmt.Sprintf(
		"🦊 **Invite Yumeko**\n%s",
		inviteURL,
	)

	respondEphemeral(s, i, content)
}
