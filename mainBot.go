package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/mcstatus-io/mcutil/v3"
)

var (
	host        = "play.ndsmc.ru"
	port uint16 = 25565
)

func mainBot(ctx context.Context, discordToken string) {
	if discordToken == "" {
		log.Println("Ошибка: Не установлен токен Discord бота.")
		return
	}

	dg, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		log.Printf("Ошибка при создании сессии Discord: %s\n", err)
		return
	}

	dg.AddHandler(func(s *discordgo.Session, event *discordgo.Ready) {
		onReady(ctx, s, event) // передаем контекст в функцию onReady
	})

	err = dg.Open()
	if err != nil {
		log.Printf("Ошибка при открытии соединения: %s\n", err)
		return
	}

	defer dg.Close()

	log.Println("Бот успешно запущен. Ожидание событий...")

	// Ожидаем отмены контекста
	<-ctx.Done()
}

func onReady(ctx context.Context, s *discordgo.Session, event *discordgo.Ready) {
	log.Printf("Бот %s подключился к Discord!\n", event.User.Username)

	// Создаем тикер для обновления статуса каждую минуту
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			updateStatus(s)
		case <-ctx.Done():
			return
		}
	}
}

func updateStatus(s *discordgo.Session) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	statusRaw, err := mcutil.StatusRaw(ctx, host, port)
	if err != nil || statusRaw == nil {
		log.Println("Сервер недоступен")
		s.UpdateGameStatus(0, "Сервер недоступен")
		return
	}

	var activity string
	if statusRaw["players"] != nil {
		players := statusRaw["players"].(map[string]interface{})
		online := int(players["online"].(float64))
		max := int(players["max"].(float64))
		activity = fmt.Sprintf("Онлайн: %d/%d", online, max)
	} else {
		activity = "Онлайн: 0/0"
	}

	s.UpdateGameStatus(0, activity)
}
