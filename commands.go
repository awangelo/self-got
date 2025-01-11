package main

import (
	"fmt"
	"log"
	"runtime"
	execCommand "self_got/exec"
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
		Help: "Displays running time and memory usage of the selfbot" + "\n\n" +
			"Example:" + "\n\n" + "Running for 5m1.793249138s\nUsing: 3 MB",
		Exec: infoCommand,
	}.add()
	command{
		Name: "bounce",
		Help: "Generates a bouncing gif based on the given image/url",
		Exec: execCommand.Bounce,
	}.add()
	command{
		Name: "remind",
		Help: "Reminds the user after a given time" + "\n\n" +
			"Example:" + "\n" + "\\remind 5m bath the cat" + "\n\n" +
			"Response after the time:" + "\n" + "@user I'm reminding you about: bath the cat",
		Exec: execCommand.Remind,
	}.add()
	command{
		Name: "ocr",
		Help: "Performs OCR on the given image",
		Exec: execCommand.Ocr,
	}.add()
}

func infoCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	memInfo := fmt.Sprintf(
		"HeapAlloc: %v MB\n"+
			"Sys: %v MB\n"+
			"Running for: %v\n",
		mem.HeapAlloc/1024/1024,
		mem.Sys/1024/1024,
		time.Since(loginTime))

	s.ChannelMessageSend(m.ChannelID, memInfo)
}

func parseCommand(s *discordgo.Session, m *discordgo.MessageCreate, message string) {
	msglist := strings.Fields(message)
	command := msglist[0]

	log.Printf("got: %v\n", msglist)
	log.Printf("replied %v to: %v\n", command, m.Author.Username)

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
