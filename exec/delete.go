package exec

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var deleteCancelChan chan struct{}

func Delete(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if deleteCancelChan == nil {
		deleteCancelChan = make(chan struct{})
	} else {
		close(deleteCancelChan)
		deleteCancelChan = make(chan struct{})
	}

	if len(args) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Not enough arguments")
		return
	}

	deleteTime := time.Now()
	deleted := 0
	delay := 500 * time.Millisecond

	switch args[0] {
	case "stop":
		close(deleteCancelChan)
		deleteCancelChan = nil
		return

	case "all":
		deleteAllMessages(s, m, &deleted, &delay, deleteCancelChan)

	default:
		count, err := strconv.Atoi(args[0])
		if err != nil || count < 1 {
			s.ChannelMessageSend(m.ChannelID, "Invalid number of messages")
			return
		}
		if count > 100 {
			count = 100
		}
		deleteNMessages(s, m, &deleted, &delay, count, deleteCancelChan)
	}

	finalMessage(s, m, deleted, time.Since(deleteTime))
}

func deleteNMessages(s *discordgo.Session, m *discordgo.MessageCreate, deleted *int, delay *time.Duration, limit int, cancelChan chan struct{}) {
	msgs, err := s.ChannelMessages(m.ChannelID, limit, "", "", "")
	if err != nil || len(msgs) == 0 {
		return
	}

	for _, msg := range msgs {
		select {
		case <-cancelChan:
			return
		default:
		}

		err := s.ChannelMessageDelete(m.ChannelID, msg.ID)
		handleRateLimit(err, delay)
		time.Sleep(*delay)
		*deleted++
	}
}

func deleteAllMessages(s *discordgo.Session, m *discordgo.MessageCreate, deleted *int, delay *time.Duration, cancelChan chan struct{}) {
	for {
		select {
		case <-cancelChan:
			return
		default:
		}

		msgs, err := s.ChannelMessages(m.ChannelID, 10, "", "", "")
		if err != nil || len(msgs) == 0 {
			return
		}

		for _, msg := range msgs {
			select {
			case <-cancelChan:
				return
			default:
			}

			err := s.ChannelMessageDelete(m.ChannelID, msg.ID)
			handleRateLimit(err, delay)
			time.Sleep(*delay)
			*deleted++
		}
	}
}

func handleRateLimit(err error, currentDelay *time.Duration) {
	if err != nil && strings.Contains(err.Error(), "rate limited") {
		*currentDelay *= 2

		if *currentDelay > 10*time.Second {
			*currentDelay = 10 * time.Second
		}
	}

	// Slowly decrease the delay if not penalized
	*currentDelay = time.Duration(float64(*currentDelay) / 1.05)
}

func finalMessage(s *discordgo.Session, m *discordgo.MessageCreate, deleted int, duration time.Duration) {
	s.ChannelMessageSend(m.ChannelID,
		fmt.Sprintf("Deleted %d messages in %v.", deleted, duration))
}
