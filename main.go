package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Создаем контекст с отменой
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Запускаем ботов в горутинах
	go mainBot(ctx, os.Getenv("MAIN_BOT_TOKEN"))
	go photoBot(ctx, os.Getenv("PHOTO_BOT_TOKEN"), os.Getenv("PHOTO_BOT_CHANNEL_IDS"))

	// Ожидаем сигнал остановки
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	// Отменяем контекст при получении сигнала остановки
	cancel()
}
