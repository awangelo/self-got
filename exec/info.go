package exec

import (
	"fmt"
	"runtime"

	"github.com/bwmarrin/discordgo"
)

func Info(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	memInfo := fmt.Sprintf(
		"HeapAlloc: %v MB\n"+
			"Sys: %v MB\n",
		mem.HeapAlloc/1024/1024,
		mem.Sys/1024/1024,
	)
	s.ChannelMessageSend(m.ChannelID, memInfo)
}
