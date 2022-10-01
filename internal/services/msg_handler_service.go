package services

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/model"
)

type MessageSender interface {
	SendMessage(text string, userID int64) error
}

type Storage interface {
	Save(*model.Spending)
	GetStatsBy(model.ReportType) (time.Time, time.Time, map[model.Category]decimal.Decimal)
}

type MessageHandlerService struct {
	tgClient MessageSender
	storage  Storage
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

var dtTemplate = "02-01-2006"

func NewMessageHandlerService(tgClient MessageSender, storage Storage) *MessageHandlerService {
	return &MessageHandlerService{
		tgClient: tgClient,
		storage:  storage,
	}
}

func (s *MessageHandlerService) HandleMsg(msg *model.Message) error {
	tokens := strings.Split(msg.Text, " ")
	if len(tokens) == 0 {
		return nil
	}

	switch tokens[0] {
	case "/start":
		return s.tgClient.SendMessage("hello", msg.UserID)
	case "/help":
		return s.tgClient.SendMessage(helpMsg, msg.UserID)
	case "/add":
		resp := ""
		if len(tokens) == 4 {
			catStr := tokens[1]
			sumStr := tokens[2]
			dtStr := tokens[3]
			if cat, err := strconv.Atoi(catStr); err != nil {
				resp = "category must be a number"
			} else if sum, err := decimal.NewFromString(sumStr); err != nil {
				resp = "sum  must be a number"
			} else if dt, err := time.Parse(dtTemplate, dtStr); err != nil {
				resp = "wrong date format"
			} else if cat > int(model.Other) || cat < 0 {
				resp = "wrong category"
			} else {
				s.storage.Save(model.NewSpending(sum, model.Category(cat), dt))
				resp = "added"
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
				resp = formatStats(s.storage.GetStatsBy(model.Week))
			case "m":
				resp = formatStats(s.storage.GetStatsBy(model.Month))
			case "y":
				resp = formatStats(s.storage.GetStatsBy(model.Year))
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

func formatStats(start time.Time, end time.Time, r map[model.Category]decimal.Decimal) string {
	if len(r) != 0 {
		result := fmt.Sprintf("from: %v, to: %v\n", start.Format(dtTemplate), end.Format(dtTemplate))
		cats := make([]model.Category, 0, len(r))
		for k := range r {
			cats = append(cats, k)
		}
		sort.Slice(cats, func(i, j int) bool { return int(cats[i]) < int(cats[j]) })
		for i := 0; i < len(cats); i++ {
			result += fmt.Sprintf("%v - %v\n", cats[i], r[cats[i]])
		}
		return result
	}
	return "no data"
}
