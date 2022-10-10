package tg

import (
	"context"
	"log"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/model"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/services"
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

		log.Println("listening for messages")

		for {
			select {
			case <-ctx.Done():
				c.client.StopReceivingUpdates()
				log.Println("stop listening messages")
				return
			case update := <-updates:
				if update.Message != nil { // If we got a message
					log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

					err := handler.HandleMsg(&model.Message{
						Text:   update.Message.Text,
						UserID: update.Message.From.ID,
					})
					if err != nil {
						log.Println("error processing message:", err)
					}
				}
			}
		}
	})
}
