package main

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	channelIDs []string
)

func photoBot(ctx context.Context, discordToken string, channelIDsRaw string) {
	channelIDs = strings.Split(channelIDsRaw, ",")
	if discordToken == "" {
		log.Println("Ошибка: Не установлен токен Discord бота.")
		return
	}

	dg, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		log.Fatalf("Ошибка при создании сессии Discord: %s", err)
		return
	}

	dg.AddHandler(onMessage)

	err = dg.Open()
	if err != nil {
		log.Fatalf("Ошибка при открытии соединения: %s", err)
		return
	}

	log.Println("Бот успешно запущен. Ожидание событий...")

	// Ожидаем отмены контекста
	<-ctx.Done()
}

func onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	if !contains(channelIDs, m.ChannelID) {
		return
	}

	currentDate := time.Now()
	now := currentDate.Format("02.01.2006, 15:04")
	username := m.Author.Username
	messageID := m.ID

	if len(m.Attachments) >= 1 {
		s.MessageReactionAdd(m.ChannelID, m.ID, "👍")
		m.ChannelID = createThread(s, m.ChannelID, messageID, now, username).ID
	} else {
		s.ChannelMessageDelete(m.ChannelID, m.ID)
	}
}

func createThread(s *discordgo.Session, channelID, messageID, now, username string) *discordgo.Channel {
	thread, err := s.MessageThreadStartComplex(channelID, messageID, &discordgo.ThreadStart{
		Name:                now + " | " + username,
		AutoArchiveDuration: 60,
		Invitable:           false,
		RateLimitPerUser:    10,
	})
	if err != nil {
		log.Printf("Ошибка при создании треда: %s", err)
	}
	return thread
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
