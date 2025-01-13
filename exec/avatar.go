package exec

import (
	"github.com/bwmarrin/discordgo"
)

func Avatar(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) > 1 {
		s.ChannelMessageSend(m.ChannelID, "Too many arguments")
		return
	}

	var targetUser *discordgo.User
	if m.ReferencedMessage != nil {
		targetUser = m.ReferencedMessage.Author
		icon := targetUser.AvatarURL("2048")
		s.ChannelMessageSend(m.ChannelID, icon)
	} else {
		taget := m.Mentions[0]

		icon := taget.AvatarURL("2048")
		s.ChannelMessageSend(m.ChannelID, icon)
	}
}
