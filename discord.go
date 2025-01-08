package main

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func testToken() error {
	fmt.Println("trying to login...")

	dg, err := discordgo.New(cfg.Token)
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
	defer dg.Close()
}

func handleReady(s *discordgo.Session, m *discordgo.Ready) {
	fmt.Printf("log-in successful!\nlog-in time: %.2f\n", time.Since(loginTime).Seconds())
	fmt.Printf("Joined %d guilds\n", len(m.Guilds))
	fmt.Printf("m.PrivateChannels: %v\n", m.PrivateChannels)
}

func handleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "pong")
	}
}
