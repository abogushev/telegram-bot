package tg

import (
	"context"
	"github.com/opentracing/opentracing-go/ext"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	. "gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/logger"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/model"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/observability"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/services"
	"go.uber.org/zap"
)

type Client struct {
	client  *tgbotapi.BotAPI
	runOnce sync.Once
}

func New(token string) (*Client, error) {
	client, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, errors.Wrap(err, "NewBotAPI")
	}

	return &Client{
		client:  client,
		runOnce: sync.Once{},
	}, nil
}

func (c *Client) SendMessage(text string, userID int64) error {
	_, err := c.client.Send(tgbotapi.NewMessage(userID, text))
	if err != nil {
		return errors.Wrap(err, "client.Send")
	}
	return nil
}

func (c *Client) ListenUpdates(handler *services.MessageHandlerService, ctx context.Context) {
	c.runOnce.Do(func() {
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60

		updates := c.client.GetUpdatesChan(u)

		Log.Info("listening for messages")

		for {
			select {
			case <-ctx.Done():
				c.client.StopReceivingUpdates()
				Log.Info("stop listening messages")
				return
			case update := <-updates:

				if update.Message != nil { // If we got a message
					Log.Info("inocming msg", zap.String("username", update.Message.From.UserName), zap.String("text", update.Message.Text))

					span, newCtx := opentracing.StartSpanFromContext(ctx, "handling message")

					observability.LogRequest(func() error {
						err := handler.HandleMsg(&model.Message{
							Text:   update.Message.Text,
							UserID: update.Message.From.ID,
						}, newCtx)
						if err != nil {
							Log.Error("error processing message:", zap.Error(err))
							ext.Error.Set(span, true)
						}
						return err
					})

					span.Finish()
				}
			}
		}
	})
}
