package main

import (
	"strings"

	execCommand "github.com/awangelo/self-got/exec"

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
			"Example:" + "\n\n" + "HeapAlloc: 18 MB \nSys: 53 MB",
		Exec: execCommand.Info,
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
	command{
		Name: "delete",
		Help: "Deletes the given number of messages" + "\n\n" +
			"Example:" + "\n" + "\\delete 5" + "\n\n" +
			"Example:" + "\n" + "\\delete all" + "\n\n" +
			"The deletion will stop if the bot receives the command 'delete stop'",
		Exec: execCommand.Delete,
	}.add()
	command{
		Name: "icon",
		Help: "Displays the icon of the given user" + "\n\n" +
			"Example:" + "\n" + "\\icon @user",
		Exec: execCommand.Icon,
	}.add()
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
