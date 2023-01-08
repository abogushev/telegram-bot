package main

import (
	"context"
	"fmt"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/model"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/cache"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/clients/tg"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/config"
	. "gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/logger"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/observability"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/services"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/storage"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/storage/pgdatabase"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/storage/pgdatabase/migrations"
	"go.uber.org/zap"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	http.Handle("/metrics", promhttp.Handler())

	observability.InitTracing(Log, "tg-bot")

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

	redisClient, err := cache.NewRedisCache(cfg.CacheHost, cfg.CachePort)
	if err != nil {
		Log.Fatal("redis init failed:", zap.Error(err))
	}
	Log.Info("starting up redis client...")

	spendigStorage := storage.NewCachedSpendingStorage(pgdatabase.NewSpendingStorage(ctx, db), redisClient)
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

	reportProducer, err := services.NewReportProducer(ctx, cfg)
	if err != nil {
		Log.Fatal("reportProducer init failed", zap.Error(err))
	}
	Log.Info("init reportProducer")
	reportResultCh := make(chan *model.Report, 10)

	if err := services.RunGRPCServer(ctx, reportResultCh); err != nil {
		Log.Fatal("RunGRPCServer failed", zap.Error(err))
	}
	Log.Info("init RunGRPCServer")

	handler := services.NewMessageHandlerService(
		tgClient,
		spendingService,
		currencyService,
		categoryService,
		stateService,
		reportProducer,
		reportResultCh,
	)
	Log.Info("init msg handler")

	go tgClient.ListenUpdates(handler, ctx)

	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%d", 8080), nil)
		if err != nil {
			Log.Fatal("error starting http server", zap.Error(err))
		}
	}()

	<-ctx.Done()
	Log.Info("gracefull shutdown...)")
	<-time.NewTimer(cfg.GracefullShutdownTimeout).C
}
