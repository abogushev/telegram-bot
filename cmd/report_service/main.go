package main

import (
	"context"
	. "gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/logger"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/report_service"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/storage/pgdatabase"
	"go.uber.org/zap"
)

func main() {
	db, err := pgdatabase.InitDB(context.Background(), "user=postgres password=postgres host=localhost dbname=postgres sslmode=disable")
	if err != nil {
		Log.Fatal("db init failed:", zap.Error(err))
	}
	storage := report_service.ReportStorage{DB: db}
	sender := report_service.NewReportResultSender(context.Background())
	if err := report_service.RunReportConsumer(context.Background(), sender, &storage); err != nil {
		Log.Fatal("failed to run consumer", zap.Error(err))
	}
	Log.Info("report service started successfully")
	<-context.Background().Done()
}
