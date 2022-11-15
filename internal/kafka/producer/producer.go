package producer

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	. "gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/logger"
	"go.uber.org/zap"
)

type Producer struct {
	p     sarama.AsyncProducer
	topic string
}

func NewProducer(ctx context.Context, topic string, brokers []string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_5_0_0
	// So we can know the partition and offset of messages.
	config.Producer.Return.Successes = true

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("starting Sarama producer: %w", err)
	}

	// We will log to STDOUT if we're not able to produce messages.
	go func() {
		for {
			select {
			case err := <-producer.Errors():
				Log.Error("Failed to write message:", zap.Error(err))
			case <-ctx.Done():
				err := producer.Close()
				if err != nil {
					Log.Error("failed to close producer")
					return
				}
				Log.Info("shutdown kafka producer")
				return
			}
		}
	}()

	return &Producer{producer, topic}, nil
}

func (producer *Producer) Send(key string, value string) {
	msg := sarama.ProducerMessage{
		Topic: producer.topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(value),
	}
	producer.p.Input() <- &msg
	successMsg := <-producer.p.Successes()
	Log.Info("Successful to write message", zap.Int64("offset", successMsg.Offset))
}
