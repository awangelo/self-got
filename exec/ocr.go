package exec

import (
	"io"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/otiai10/gosseract/v2"
)

func Ocr(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	var imageURL string

	if len(args) > 1 {
		s.ChannelMessageSend(m.ChannelID, "Too many arguments")
		return
	}

	// User provided an attachment
	imageURL, err := getImageFromMessage(m)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	client := gosseract.NewClient()
	defer client.Close()

	resp, err := http.Get(imageURL)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Failed to download image")
		return
	}
	defer resp.Body.Close()
	imgBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Failed to read image")
		return
	}

	if len(args) == 1 {
		client.SetLanguage(args[0])
	}

	client.SetImageFromBytes(imgBytes)
	text, _ := client.Text()

	s.ChannelMessageSend(m.ChannelID, text)
}
