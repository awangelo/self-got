package exec

import (
	"fmt"
	"net/url"
	"strings"
	"time"

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

	return "", nil
}

func getImageNameFromMessage(m *discordgo.MessageCreate) string {
	if len(m.Attachments) == 1 {
		return m.Attachments[0].Filename
	}

	repliedMessage := m.ReferencedMessage
	if repliedMessage != nil {
		if len(repliedMessage.Attachments) == 1 {
			return repliedMessage.Attachments[0].Filename
		}
	}

	// return random name
	return m.ID[:8]
}

func handleRateLimit(err error, currentDelay *time.Duration) {
	if err != nil && strings.Contains(err.Error(), "rate limited") {
		*currentDelay *= 2

		if *currentDelay > 10*time.Second {
			*currentDelay = 10 * time.Second
		}
	}

	// Slowly decrease the delay if not penalized
	*currentDelay = time.Duration(float64(*currentDelay) / 1.05)
}
