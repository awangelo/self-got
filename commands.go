package main

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type command struct {
	Name string
	Help string
	Exec func(*discordgo.Session, *discordgo.MessageCreate, []string)
}

var (
	commMap = make(map[string]command)
)

func prepareCommands() {
	command{
		Name: "info",
		Help: "Displays running time and memory usage of the selfbot" + "\n\n" + "Example:" + "\n\n" + "Running for 5m1.793249138s\nUsing: 3 MB",
		Exec: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			var mem runtime.MemStats
			runtime.ReadMemStats(&mem)
			memUsage := fmt.Sprintf("Running for %v\nUsing: %v MB\n", time.Since(loginTime), mem.Alloc/1024/1024)
			s.ChannelMessageSend(m.ChannelID, memUsage)
		},
	}.add()
	command{
		Name: "ping",
		Help: "Responds with Pong!",
		Exec: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			s.ChannelMessageSend(m.ChannelID, "Pong!")
		},
	}.add()
}

func parseCommand(s *discordgo.Session, m *discordgo.MessageCreate, message string) {
	msglist := strings.Fields(message)
	command := msglist[0]

	fmt.Printf("got: %v\n", msglist)
	fmt.Printf("replied %v to: %v\n", command, m.Author.Username)

	if command == commMap[command].Name {
		commMap[command].Exec(s, m, msglist[1:])
	}
}

func (c command) add() {
	commMap[c.Name] = c
}

func getCommandNames() []string {
	commandNames := make([]string, 0, len(commMap))

	for name := range commMap {
		commandNames = append(commandNames, name)
	}

	return commandNames
}
