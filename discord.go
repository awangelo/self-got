package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func testToken() error {
	fmt.Println("trying to login...")

	var err error
	dg, err = discordgo.New(cfg.Token)
	if err != nil {
		return err
	}

	if _, err = dg.User("@me"); err != nil {
		return err
	}

	fmt.Println("token is valid!")

	return nil
}

func connectToDiscord() {
	fmt.Println("opening websocket...")
	if err := dg.Open(); err != nil {
		fmt.Println(err)
	}
}

func handleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !strings.HasPrefix(m.Content, cfg.Prefix) {
		return
	}

	parseCommand(s, m, strings.TrimPrefix(m.Content, cfg.Prefix))
}
