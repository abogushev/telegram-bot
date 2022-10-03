package services

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
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
	tgClient                 MessageSender
	storage                  Storage
	currentCurrencyType      model.CurrencyType
	currentCurrencyTypeMutex sync.RWMutex
}

var helpMsg = `
/help - call this help
/add [category] [sum] - add spending, template
/categories - show all categories
/report [type] - show report. type: w - week, m - month, y - year
/currency [type] - change currency. type: 0 - USD, 1 - CNY, 2 - EUR, 3 - RUB
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
	} else if cat > int(model.Other) || cat < 0 {
		return "", errors.New("wrong category")
	} else {
		s.storage.Save(model.NewSpending(sum, model.Category(cat), dt))
		return "added", nil
	}
}

func handleCategories() string {
	return categoriesMsg
}

func (s *MessageHandlerService) handleReport(strs []string) (string, error) {
	switch strs[1] {
	case "w":
		return formatStats(s.storage.GetStatsBy(model.Week)), nil
	case "m":
		return formatStats(s.storage.GetStatsBy(model.Month)), nil
	case "y":
		return formatStats(s.storage.GetStatsBy(model.Year)), nil
	default:
		return "", errWrongFormat
	}
}

func (s *MessageHandlerService) handleCurrencyChange(strs []string) (string, error) {
	if i, err := strconv.Atoi(strs[1]); err != nil {
		return "", errWrongFormat
	} else if cur, err := model.ParseCurrencyType(i); err != nil {
		return "", err
	} else {
		s.currentCurrencyTypeMutex.Lock()
		s.currentCurrencyType = cur
		s.currentCurrencyTypeMutex.Unlock()
		return "successfully changed", nil
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
