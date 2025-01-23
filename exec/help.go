package exec

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func Help(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	commands := []struct {
		Name  string
		Value string
	}{
		{
			Name:  "help",
			Value: "- Displays the list of commands.",
		},
		{
			Name:  "info",
			Value: "- Displays running time and memory usage of the selfbot.",
		},
		{
			Name:  "bounce",
			Value: "- Generates a bouncing gif based on the given image/url.",
		},
		{
			Name:  "remind",
			Value: "- Reminds the user after a given time.",
		},
		{
			Name:  "ocr",
			Value: "- Performs OCR on the given image.",
		},
		{
			Name:  "delete",
			Value: "- Deletes the given number of messages.",
		},
		{
			Name:  "avatar",
			Value: "- Displays the avatar of the given user.",
		},
		{
			Name:  "nuke",
			Value: "- Deletes all channels and roles in the server.",
		},
	}

	var sb strings.Builder
	sb.WriteString("## **Available Commands:**\n")
	for _, cmd := range commands {
		sb.WriteString("**" + cmd.Name + "**: " + cmd.Value + "\n")
	}
	sb.WriteString("-# keep in mind that image commands work with image attachments, replies, and urls.")

	message := sb.String()

	s.ChannelMessageSend(m.ChannelID, message)
}
