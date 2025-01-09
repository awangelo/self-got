package exec

import (
	"io"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/otiai10/gosseract"
)

func Ocr(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if m.Embeds == nil {
		s.ChannelMessageSend(m.ChannelID, "You need to provide an image")
		return
	}
	if len(args) > 1 {
		s.ChannelMessageSend(m.ChannelID, "Too many arguments")
		return
	}

	client := gosseract.NewClient()
	defer client.Close()

	image := m.Embeds[0].Image.URL
	resp, err := http.Get(image)
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

	lang := args[0]
	if lang == "" {
		client.SetLanguage(lang)
	}

	client.SetImageFromBytes(imgBytes)
	text, _ := client.Text()

	s.ChannelMessageSend(m.ChannelID, text)
}
