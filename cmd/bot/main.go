package main

import (
	"log"

	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/clients/tg"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/config"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/services"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/storage"
)

func main() {
	config, err := config.New()
	if err != nil {
		log.Fatal("config init failed:", err)
	}

	tgClient, err := tg.New(config)
	if err != nil {
		log.Fatal("tg client init failed:", err)
	}

	handler := services.NewMessageHandlerService(tgClient, storage.NewInMemoryStorage())

	tgClient.ListenUpdates(handler)
}
