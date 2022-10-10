package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/clients/tg"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/config"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/services"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/storage"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	cfg, err := config.New()
	if err != nil {
		log.Fatal("config init failed:", err)
	}

	tgClient, err := tg.New(cfg.Token)

	if err != nil {
		log.Fatal("tg client init failed:", err)
	}

	spendigStorage := storage.NewInMemorySpendingStorage()
	currencyStorage := storage.NewCurrencyStorage()
	spendingService := services.NewSpendingService(spendigStorage, currencyStorage)
	currencyService := services.NewCurrencyService()

	currencyService.RunUpdateCurrenciesDaemon(ctx, cfg.UpdateCurrenciesInterval, currencyStorage)

	handler := services.NewMessageHandlerService(tgClient, spendingService)

	go tgClient.ListenUpdates(handler, ctx)

	<-ctx.Done()
	log.Println("gracefull shutdown...)")
	<-time.NewTimer(cfg.GracefullShutdownTimeout).C
}
