package exec

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func Remind(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Usage: remind <time> <message>")
		return
	}

	duration, err := time.ParseDuration(args[0])
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Invalid time, use a format like 10m, 1h, etc.")
		return
	}

	message := strings.Join(args[1:], " ")

	time.AfterFunc(duration, func() {
		reminderMessage := fmt.Sprintf("<@%s> I'm reminding you about: %s", m.Author.ID, message)
		s.ChannelMessageSend(m.ChannelID, reminderMessage)
	})
}
