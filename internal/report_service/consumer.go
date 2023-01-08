package report_service

import (
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	. "gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/logger"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/model"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/time_util"
	"go.uber.org/zap"
)

var (
	KafkaTopic         = "report"
	KafkaConsumerGroup = "report-consumer-group"
	BrokersList        = []string{"localhost:9092"}
)

func RunReportConsumer(ctx context.Context, reportResultService *ReportResultSender, reportStorage *ReportStorage) error {
	config := sarama.NewConfig()
	config.Version = sarama.V2_5_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	// Create consumer group
	consumerGroup, err := sarama.NewConsumerGroup(BrokersList, KafkaConsumerGroup, config)
	if err != nil {
		Log.Error("starting consumer group: %w", zap.Error(err))
		return err
	}

	err = consumerGroup.Consume(ctx, []string{KafkaTopic}, &Consumer{reportResultService, reportStorage})
	if err != nil {
		Log.Error("consuming via handler: %w", zap.Error(err))
		return err
	}
	return nil
}

// Consumer represents a Sarama consumer group consumer.
type Consumer struct {
	reportResultService *ReportResultSender
	reportStorage       *ReportStorage
}

// Setup is run at the beginning of a new session, before ConsumeClaim.
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	Log.Info("consumer - setup")
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited.
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	Log.Info("consumer - cleanup")
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		consumer.handleMessage(message)
		Log.Info("msg handled")
		session.MarkMessage(message, "")
	}
	return nil
}

func (consumer *Consumer) handleMessage(m *sarama.ConsumerMessage) {
	request := model.ReportRequest{}
	err := json.Unmarshal(m.Value, &request)
	if err != nil {
		Log.Error("faield to parse msg", zap.Error(err))
		return
	}
	start, err := time_util.DateToTime(request.Start)
	if err != nil {
		Log.Error("faield to parse start field", zap.Error(err))
		return
	}
	end, err := time_util.DateToTime(request.End)
	if err != nil {
		Log.Error("failed to parse end field", zap.Error(err))
		return
	}
	result, err := consumer.reportStorage.getStatsBy(context.Background(), start, end)
	if err != nil {
		Log.Error("failed to get stat from db")
	}
	consumer.reportResultService.Send(context.Background(), request.UserId, request.Start, request.End, result)
}
