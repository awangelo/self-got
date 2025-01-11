package exec

import (
	"fmt"
	"net/url"

	"github.com/bwmarrin/discordgo"
)

func getImageFromMessage(m *discordgo.MessageCreate) (string, error) {
	if len(m.Attachments) == 1 {
		return m.Attachments[0].URL, nil
	}

	repliedMessage := m.ReferencedMessage
	if repliedMessage != nil {
		if len(repliedMessage.Attachments) == 1 {
			return repliedMessage.Attachments[0].URL, nil
		}

		_, err := url.ParseRequestURI(repliedMessage.Content)
		if err != nil {
			return "", fmt.Errorf("Invalid URL provided")
		}

		return repliedMessage.Content, nil
	}

	return "", fmt.Errorf("You need to provide a valid image or URL")
}
