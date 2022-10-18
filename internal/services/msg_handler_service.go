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
	Save(model.Spending) error
	GetStatsBy(time.Time, time.Time) (map[string]decimal.Decimal, string, error)
}

type CurrencyService interface {
	UpdateCurrentCurrency(c string) error
	GetAll() []model.Currency
}

type CategoryService interface {
	GetAll() []model.Category
}

type MessageHandlerService struct {
	tgClient        MessageSender
	spendingService SpendingService
	currencyService CurrencyService
	categoryService CategoryService
}

var helpMsg = `
/help - call this help
/categories - show all categories
/currencies - show all currencies
/add [category] [sum] - add spending
/report [type] - show report. type: w - week, m - month, y - year
/currency [type] - change currency
`

var dtTemplate = "02-01-2006"

func NewMessageHandlerService(
	tgClient MessageSender,
	spendingService SpendingService,
	currencyService CurrencyService,
	categoryService CategoryService) *MessageHandlerService {

	return &MessageHandlerService{
		tgClient:        tgClient,
		spendingService: spendingService,
		currencyService: currencyService,
		categoryService: categoryService,
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
		resp = "hello"

	case "/help":
		resp = helpMsg

	case "/add":
		resp = handleF(tokens, 4, s.handleAdd)

	case "/categories":
		resp = s.handleCategories()

	case "/report":
		resp = handleF(tokens, 2, s.handleReport)

	case "/currencies":
		resp = s.handleCurrencies()

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
func genListMsg(els []string) string {
	var sb strings.Builder
	for i := 0; i < len(els); i++ {
		sb.WriteString("- ")
		sb.WriteString(els[i])
		sb.WriteString("\n")
	}
	return sb.String()
}
func (s *MessageHandlerService) handleCurrencies() string {
	allCrns := s.currencyService.GetAll()
	els := make([]string, len(allCrns))

	for i := 0; i < len(allCrns); i++ {
		els[i] = allCrns[i].Code
	}

	return genListMsg(els)
}

func (s *MessageHandlerService) handleCategories() string {
	allCats := s.categoryService.GetAll()
	els := make([]string, len(allCats))

	for i := 0; i < len(allCats); i++ {
		els[i] = strconv.Itoa(allCats[i].Id)
	}

	return genListMsg(els)
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
	} else if err := s.spendingService.Save(model.NewSpending(sum, cat, dt)); err != nil {
		return "", err
	}
	return "added", nil
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
	if err := s.currencyService.UpdateCurrentCurrency(strs[1]); err != nil {
		return "", err
	}
	return "successfully changed", nil
}

func formatStats(start time.Time, end time.Time, r map[string]decimal.Decimal, currencyCode string) string {
	if len(r) != 0 {
		result := fmt.Sprintf("from: %v, to: %v\n", start.Format(dtTemplate), end.Format(dtTemplate))
		cats := make([]string, 0, len(r))
		for k := range r {
			cats = append(cats, k)
		}
		sort.Slice(cats, func(i, j int) bool { return cats[i] < cats[j] })
		for i := 0; i < len(cats); i++ {
			result += fmt.Sprintf("%v - %v %v\n", cats[i], r[cats[i]].Round(2), currencyCode)
		}
		return result
	}
	return "no data"
}
