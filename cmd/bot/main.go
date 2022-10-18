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
	log.Println("init cnfg")

	tgClient, err := tg.New(cfg.Token)
	if err != nil {
		log.Fatal("tg client init failed:", err)
	}
	log.Println("init tgClient")


	db, err := pgdatabase.InitDB(ctx, "user=postgres password=postgres dbname=postgres sslmode=disable")
	migrations.Up("postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable", "internal/storage/pgdatabase/migrations")

	if err != nil {
		log.Fatal("db init failed:", err)
	}
	log.Println("init db")

	spendigStorage := pgdatabase.NewSpendingStorage(ctx, db)
	log.Println("init spendigStorage")
	currencyStorage := pgdatabase.NewCurrencyStorage(ctx, db)
	log.Println("init currencyStorage")
	currencyService, err := services.NewCurrencyService(currencyStorage)
	if err != nil {
		log.Fatal("currencyService init failed", err)
	}
	log.Println("init currencyService")
	currencyService.RunUpdateCurrenciesDaemon(ctx, cfg.UpdateCurrenciesInterval)
	log.Println("run RunUpdateCurrenciesDaemon")

	spendingService := services.NewSpendingService(spendigStorage, currencyService)
	log.Println("init spendingService")

	categoryStorage := pgdatabase.NewCategoryStorage(ctx, db)
	log.Println("init categoryStorage")

	categoryService, err := services.NewCategoryService(categoryStorage)
	if err != nil {
		log.Fatal("categoryService init failed", err)
	}
	log.Println("init categoryService")

	handler := services.NewMessageHandlerService(tgClient, spendingService, currencyService, categoryService)
	log.Println("init msg handler")
	
	go tgClient.ListenUpdates(handler, ctx)

	<-ctx.Done()
	log.Println("gracefull shutdown...)")
	<-time.NewTimer(cfg.GracefullShutdownTimeout).C
}
