package report_service

import (
	"context"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/api"
	. "gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type ReportResultSender struct {
	client api.ReportClient
}

func NewReportResultSender(ctx context.Context) *ReportResultSender {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := api.NewReportClient(conn)
	go func() {
		<-ctx.Done()
		if err := conn.Close(); err != nil {
			Log.Error("failed to close conn", zap.Error(err))
		}
	}()
	Log.Info("ReportResultSender created")
	return &ReportResultSender{c}
}

func (s *ReportResultSender) Send(ctx context.Context, userId int64, start string, end string, data map[string]float64) {
	_, err := s.client.Send(ctx, &api.ReportResult{UserId: userId, Start: start, End: end, Data: data})
	if err != nil {
		Log.Error("failed on send request", zap.Error(err))
		return
	}
	Log.Info("request sent successfully")
}
