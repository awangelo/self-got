package exec

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

func Nuke(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) == 0 {
		s.ChannelMessageSend(m.ChannelID, "You need to confirm the server name or ID")
		return
	}

	target, err := s.GuildChannels(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error fetching server information")
		return
	}

	guild, err := s.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error fetching guild information")
		return
	}

	if len(args) != 1 || (args[0] != guild.Name && args[0] != m.GuildID) {
		s.ChannelMessageSend(m.ChannelID, "Invalid confirmation. Please provide the correct server name or ID.")
		return
	}

	delay := 500 * time.Millisecond

	go func() {
		for _, role := range guild.Roles {
			err := s.GuildRoleDelete(m.GuildID, role.ID)
			if err != nil {
				handleRateLimit(err, &delay)
			}
		}
	}()

	for _, channel := range target {
		_, err := s.ChannelDelete(channel.ID)
		if err != nil {
			handleRateLimit(err, &delay)
		}
	}
}
