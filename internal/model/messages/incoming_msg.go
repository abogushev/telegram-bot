package messages

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/storage"
)

type MessageSender interface {
	SendMessage(text string, userID int64) error
}

type Model struct {
	tgClient MessageSender
}

func New(tgClient MessageSender) *Model {
	return &Model{
		tgClient: tgClient,
	}
}

type Message struct {
	Text   string
	UserID int64
}

var helpMsg = `
/help - call this help
/add [category] [sum] - add spending, template
/categories - show all categories
/report [type] - show report. type: w - week, m - month, y - year
`

var categoriesMsg = `
categories: 
	0 - food
	1 - other 
`

func (s *Model) IncomingMessage(msg Message) error {
	tokens := strings.Split(msg.Text, " ")
	if len(tokens) == 0 {
		return nil
	}

	switch tokens[0] {
	case "/help":
		return s.tgClient.SendMessage(helpMsg, msg.UserID)
	case "/add":
		resp := ""
		if len(tokens) == 3 {
			catStr := tokens[1]
			sumStr := tokens[2]
			if cat, err := strconv.Atoi(catStr); err != nil {
				resp = "category must be a number"
			} else if sum, err := strconv.ParseFloat(sumStr, 32); err != nil {
				resp = "sum  must be a number"
			} else {
				if cat > int(storage.Other) || cat < 0 {
					resp = "wrong category"
				} else {
					storage.New(float32(sum), storage.Category(cat)).Save()
					resp = "added"
				}
			}
		} else {
			resp = "wrong format, must be: /add [category] [sum]"
		}
		return s.tgClient.SendMessage(resp, msg.UserID)

	case "/categories":
		return s.tgClient.SendMessage(categoriesMsg, msg.UserID)

	case "/report":
		resp := ""
		if len(tokens) == 2 {
			switch tokens[1] {
			case "w":
				resp = formatStats(storage.GetStatsBy(storage.Week))
			case "m":
				resp = formatStats(storage.GetStatsBy(storage.Month))
			case "y":
				resp = formatStats(storage.GetStatsBy(storage.Year))
			default:
				resp = "wrong type, must be one of w - week, m - month, y - year"
			}
		} else {
			resp = "wrong format, must be: /report [type]"
		}

		return s.tgClient.SendMessage(resp, msg.UserID)

	default:
		return s.tgClient.SendMessage("не знаю эту команду", msg.UserID)
	}
}

func formatStats(start time.Time, end time.Time, r map[storage.Category]float32) string {
	if len(r) != 0 {
		template := "02-01-2006"
		result := fmt.Sprintf("from: %v, to: %v\n", start.Format(template), end.Format(template))
		for cat, sum := range r {
			result += fmt.Sprintf("%v - %v\n", cat, sum)
		}
		return result
	}
	return "no data"
}
