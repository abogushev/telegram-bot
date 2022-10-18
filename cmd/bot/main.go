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
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/storage/pgdatabase"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/storage/pgdatabase/migrations"
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

	//	spendigStorage := storage.NewInMemorySpendingStorage()
	//	currencyStorage := storage.NewCurrencyStorage()
	db, err := pgdatabase.InitDB(ctx, "user=postgres password=postgres dbname=postgres sslmode=disable")
	migrations.Up("postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable", "internal/storage/pgdatabase/migrations")

	if err != nil {
		log.Fatal("db init failed:", err)
	}

	spendigStorage := pgdatabase.NewSpendingStorage(ctx, db)
	currencyStorage := pgdatabase.NewCurrencyStorage(ctx, db)

	currencyService, err := services.NewCurrencyService(currencyStorage)
	if err != nil {
		log.Fatal("currencyService init failed", err)
	}
	currencyService.RunUpdateCurrenciesDaemon(ctx, cfg.UpdateCurrenciesInterval)

	spendingService := services.NewSpendingService(spendigStorage, currencyService)

	categoryStorage := pgdatabase.NewCategoryStorage(ctx, db)
	categoryService, err := services.NewCategoryService(categoryStorage)
	if err != nil {
		log.Fatal("categoryService init failed", err)
	}

	handler := services.NewMessageHandlerService(tgClient, spendingService, currencyService, categoryService)

	go tgClient.ListenUpdates(handler, ctx)

	<-ctx.Done()
	log.Println("gracefull shutdown...)")
	<-time.NewTimer(cfg.GracefullShutdownTimeout).C
}
