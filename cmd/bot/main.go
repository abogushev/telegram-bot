package main

import (
	"context"
	"go.uber.org/zap"
	"os/signal"
	"syscall"
	"time"
	. "gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/logger"
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
		Log.Fatal("config init failed:%v", zap.Error(err))
	}
	Log.Info("init cnfg")

	tgClient, err := tg.New(cfg.Token)
	if err != nil {
		Log.Fatal("tg client init failed:", zap.Error(err))
	}
	Log.Info("init tgClient")

	db, err := pgdatabase.InitDB(ctx, "user=postgres password=postgres host=localhost dbname=postgres sslmode=disable")
	if err != nil {
		Log.Fatal("db init failed:", zap.Error(err))
	}
	Log.Info("starting up migrations...")

	migrations.Up("postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable", "internal/storage/pgdatabase/migrations")
	Log.Info("successfully migrated")

	Log.Info("init db")

	spendigStorage := pgdatabase.NewSpendingStorage(ctx, db)
	Log.Info("init spendigStorage")
	currencyStorage := pgdatabase.NewCurrencyStorage(ctx, db)
	Log.Info("init currencyStorage")
	currencyService, err := services.NewCurrencyService(currencyStorage)
	if err != nil {
		Log.Fatal("currencyService init failed", zap.Error(err))
	}
	Log.Info("init currencyService")
	currencyService.RunUpdateCurrenciesDaemon(ctx, cfg.UpdateCurrenciesInterval)
	Log.Info("run RunUpdateCurrenciesDaemon")

	categoryStorage := pgdatabase.NewCategoryStorage(ctx, db)
	Log.Info("init categoryStorage")

	categoryService, err := services.NewCategoryService(categoryStorage)
	if err != nil {
		Log.Fatal("categoryService init failed", zap.Error(err))
	}
	Log.Info("init categoryService")

	stateStorage := pgdatabase.NewStateStorage(ctx, db)
	stateService, err := services.NewStateService(stateStorage, ctx)
	if err != nil {
		Log.Fatal("stateService init failed", zap.Error(err))
	}
	Log.Info("init stateService")

	spendingService := services.NewSpendingService(spendigStorage, currencyService, stateService)
	Log.Info("init spendingService")

	handler := services.NewMessageHandlerService(tgClient, spendingService, currencyService, categoryService, stateService)
	Log.Info("init msg handler")

	go tgClient.ListenUpdates(handler, ctx)

	<-ctx.Done()
	Log.Info("gracefull shutdown...)")
	<-time.NewTimer(cfg.GracefullShutdownTimeout).C
}
