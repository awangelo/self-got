package exec

import (
	"os/exec"
	"runtime"

	"github.com/bwmarrin/discordgo"
)

type urlEntry struct {
	Name string
	URL  string
}

var urls = []urlEntry{
	{Name: "lens", URL: "https://lens.google.com/uploadbyurl?url="},
	{Name: "tineye", URL: "https://tineye.com/search?url="},
	{Name: "yandex", URL: "https://yandex.com/images/search?url="},
	{Name: "anime", URL: "https://trace.moe/?url="},
	{Name: "sauce", URL: "https://saucenao.com/search.php?url="},
}

func Reverse(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) != 0 {
		s.ChannelMessageSend(m.ChannelID, "Too many arguments")
		return
	}

	image, err := getImageFromMessage(m)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	if image == "" {
		s.ChannelMessageSend(m.ChannelID, "No image found")
		return
	}

	for _, entry := range urls {
		go func(entry urlEntry) {
			var cmd *exec.Cmd
			finalUrl := entry.URL + image

			if entry.Name == "yandex" {
				finalUrl += "&rpt=imageview"
			}

			if runtime.GOOS == "windows" {
				cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", finalUrl)
			} else {
				cmd = exec.Command("xdg-open", finalUrl)
			}
			err := cmd.Start()
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Failed to open URL: "+err.Error())
			}
		}(entry)
	}
}
