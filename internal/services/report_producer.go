package services

import (
	"context"
	"encoding/json"
	"fmt"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/config"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/kafka/producer"
	. "gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/logger"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/model"
	"go.uber.org/zap"
)

type ReportProducer struct {
	producer *producer.Producer
}

func NewReportProducer(ctx context.Context, config *config.Config) (*ReportProducer, error) {
	p, err := producer.NewProducer(ctx, config.TopicReport, config.KafkaBrokers)
	if err != nil {
		return nil, err
	}
	return &ReportProducer{p}, nil
}

func (p *ReportProducer) Send(request *model.ReportRequest) error {
	js, err := json.Marshal(request)
	if err != nil {
		Log.Error("failed to send report request")
		return err
	}
	p.producer.Send(fmt.Sprint(request.UserId), string(js))
	Log.Info("send report request", zap.Int64("userId", request.UserId))
	return nil
}
