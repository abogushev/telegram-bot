package services

import (
	"errors"
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

type SpendingService interface {
	Save(*model.Spending) error
	GetStatsBy(time.Time, time.Time) (map[model.Category]decimal.Decimal, model.CurrencyType, error)
	UpdateCurrentType(model.CurrencyType) error
}

type MessageHandlerService struct {
	tgClient        MessageSender
	spendingService SpendingService
}

var helpMsg = `
/help - call this help
/add [category] [sum] - add spending, template
/categories - show all categories
/report [type] - show report. type: w - week, m - month, y - year
/currency [type] - change currency. type: 1 - USD, 2 - CNY, 3 - EUR, 4 - RUB
`

var categoriesMsg = `
categories: 
	0 - food
	1 - other 
`

var dtTemplate = "02-01-2006"

func NewMessageHandlerService(tgClient MessageSender, spendingService SpendingService) *MessageHandlerService {
	return &MessageHandlerService{
		tgClient:        tgClient,
		spendingService: spendingService,
	}
}

func (s *MessageHandlerService) HandleMsg(msg *model.Message) error {
	tokens := strings.Split(msg.Text, " ")
	if len(tokens) == 0 {
		return nil
	}
	resp := ""

	switch tokens[0] {
	case "/start":
		resp = handleStart()
	case "/help":
		resp = handleHelp()
	case "/add":
		resp = handleF(tokens, 4, s.handleAdd)

	case "/categories":
		resp = handleCategories()

	case "/report":
		resp = handleF(tokens, 2, s.handleReport)

	case "/currency":
		resp = handleF(tokens, 2, s.handleCurrencyChange)

	default:
		resp = "не знаю эту команду"
	}
	return s.tgClient.SendMessage(resp, msg.UserID)
}

var errWrongFormat = errors.New("wrong format")

func handleF(strs []string, count int, handler func([]string) (string, error)) string {
	if count != len(strs) {
		return errWrongFormat.Error()
	}

	if r, err := handler(strs); err != nil {
		return err.Error()
	} else {
		return r
	}
}

func handleStart() string {
	return "hello"
}

func handleHelp() string {
	return helpMsg
}

func (s *MessageHandlerService) handleAdd(tokens []string) (string, error) {
	catStr := tokens[1]
	sumStr := tokens[2]
	dtStr := tokens[3]

	if cat, err := strconv.Atoi(catStr); err != nil {
		return "", errors.New("category must be a number")
	} else if sum, err := decimal.NewFromString(sumStr); err != nil {
		return "", errors.New("sum  must be a number")
	} else if dt, err := time.Parse(dtTemplate, dtStr); err != nil {
		return "", errors.New("wrong date format")
	} else if cat > int(model.Food) || cat < int(model.Other) {
		return "", errors.New("wrong category")
	} else if err := s.spendingService.Save(model.NewSpending(sum, model.Category(cat), dt)); err != nil {
		return "", err
	}
	return "added", nil
}

func handleCategories() string {
	return categoriesMsg
}

func (s *MessageHandlerService) handleReport(strs []string) (string, error) {
	endAt := time.Now().Truncate(24 * time.Hour)
	var startAt time.Time

	switch strs[1] {
	case "w":
		startAt = endAt.AddDate(0, 0, -7)
	case "m":
		startAt = endAt.AddDate(0, -1, 0)
	case "y":
		startAt = endAt.AddDate(-1, 0, 0)
	default:
		return "", errWrongFormat
	}

	if data, c, err := s.spendingService.GetStatsBy(startAt, endAt); err != nil {
		return "", err
	} else {
		return formatStats(startAt, endAt, data, c), nil
	}
}

func (s *MessageHandlerService) handleCurrencyChange(strs []string) (string, error) {
	if i, err := strconv.Atoi(strs[1]); err != nil {
		return "", errWrongFormat
	} else if cur, err := model.ParseCurrencyType(i); err != nil {
		return "", err
	} else if err := s.spendingService.UpdateCurrentType(cur); err != nil {
		return "failed to update", err
	}
	return "successfully changed", nil
}

func formatStats(start time.Time, end time.Time, r map[model.Category]decimal.Decimal, currency model.CurrencyType) string {
	if len(r) != 0 {
		result := fmt.Sprintf("from: %v, to: %v\n", start.Format(dtTemplate), end.Format(dtTemplate))
		cats := make([]model.Category, 0, len(r))
		for k := range r {
			cats = append(cats, k)
		}
		sort.Slice(cats, func(i, j int) bool { return int(cats[i]) < int(cats[j]) })
		for i := 0; i < len(cats); i++ {
			result += fmt.Sprintf("%v - %v %v\n", cats[i], r[cats[i]].Round(2), currency)
		}
		return result
	}
	return "no data"
}
