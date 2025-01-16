package exec

import (
	"os/exec"
	"runtime"

	"github.com/bwmarrin/discordgo"
)

var urls = map[string]string{
	"lens":   "https://lens.google.com/uploadbyurl?url=",
	"tineye": "https://tineye.com/search?url=",
	"yandex": "https://yandex.com/images/search?url=",
	"anime":  "https://trace.moe/?url=",
	"sauce":  "https://saucenao.com/search.php?url=",
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

	for n, url := range urls {
		go func(url string) {
			var cmd *exec.Cmd
			var finalUrl = url + image

			if n == "yandex" {
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
		}(url)
	}
}
