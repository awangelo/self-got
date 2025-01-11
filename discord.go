package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func createSession(cfg *config) (*discordgo.Session, error) {
	fmt.Println("trying to login...")

	dg, err := discordgo.New(cfg.Token)
	if err != nil {
		return nil, err
	}

	if _, err = dg.User("@me"); err != nil {
		return nil, err
	}

	fmt.Println("token is valid!")
	return dg, nil
}

func connectToWS(dg *discordgo.Session) error {
	fmt.Println("opening websocket...")
	if err := dg.Open(); err != nil {
		return fmt.Errorf("failed to open connection: %v", err)
	}
	return nil
}

func handleMessageCreate(cfg *config, s *discordgo.Session, m *discordgo.MessageCreate) {
	if !strings.HasPrefix(m.Content, cfg.Prefix) {
		return
	}

	parseCommand(s, m, strings.TrimPrefix(m.Content, cfg.Prefix))
}
