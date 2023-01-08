package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/shopspring/decimal"
	. "gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/logger"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/model"
	"go.uber.org/zap"
	"sort"
	"strconv"
	"strings"
	"time"
)

type MessageSender interface {
	SendMessage(text string, userID int64) error
}

type SpendingServiceI interface {
	SaveTx(context.Context, model.Spending) (decimal.Decimal, error)
	GetStatsBy(context.Context, time.Time, time.Time) (map[string]decimal.Decimal, string, error)
}

type CurrencyService interface {
	UpdateCurrentCurrency(c string) error
	GetAll() []model.Currency
}

type CategoryService interface {
	GetAll() []model.Category
}
type StateService interface {
	GetBalance() (decimal.Decimal, error)
}
type MessageHandlerService struct {
	tgClient        MessageSender
	spendingService SpendingServiceI
	currencyService CurrencyService
	categoryService CategoryService
	stateService    StateService
	reportProducer  *ReportProducer
}

var helpMsg = `
/help - call this help
/categories - show all categories
/currencies - show all currencies
/add [category] [sum] [date] - add spending 
/report [type] - show report. type: w - week, m - month, y - year
/currency [type] - change currency
`

var dtTemplate = "02-01-2006"

func NewMessageHandlerService(
	tgClient MessageSender,
	spendingService SpendingServiceI,
	currencyService CurrencyService,
	categoryService CategoryService,
	stateService StateService,
	reportProducer *ReportProducer,
	reportResultCh <-chan *model.Report) *MessageHandlerService {
	s := &MessageHandlerService{
		tgClient:        tgClient,
		spendingService: spendingService,
		currencyService: currencyService,
		categoryService: categoryService,
		stateService:    stateService,
		reportProducer:  reportProducer,
	}
	go s.reportResultListen(reportResultCh)
	return s
}

func (s *MessageHandlerService) HandleMsg(msg *model.Message, ctx context.Context) error {
	span, spanCtx := opentracing.StartSpanFromContext(ctx, "")

	defer span.Finish()

	tokens := strings.Split(msg.Text, " ")
	if len(tokens) == 0 {
		return nil
	}
	resp := ""

	switch tokens[0] {
	case "/start":
		resp = "hello"
		span.SetOperationName("msg_handler: handle cmd `/start`")
	case "/help":
		resp = helpMsg
		span.SetOperationName("msg_handler: handle cmd `/help`")
	case "/add":
		resp = handleF(span, spanCtx, tokens, 4, s.handleAdd)
		span.SetOperationName("msg_handler: handle cmd `/add`")
	case "/categories":
		resp = s.handleCategories()
		span.SetOperationName("msg_handler: handle cmd `/categories`")
	case "/report":
		//resp = handleF(span, spanCtx, tokens, 2, s.handleReport)
		resp = handleF(span, spanCtx, tokens, 2, func(ctx context.Context, i []string) (string, error) {
			r, err := s.handleReportAsync(spanCtx, msg.UserID, tokens)
			if err != nil {
				return "", err
			}
			return r, nil
		})
		span.SetOperationName("msg_handler: handle cmd `/report`")
	case "/currencies":
		resp = s.handleCurrencies()
		span.SetOperationName("msg_handler: handle cmd `/currencies`")
	case "/currency":
		resp = handleF(span, spanCtx, tokens, 2, s.handleCurrencyChange)
		span.SetOperationName("msg_handler: handle cmd `/currency`")
	case "/balance":
		resp = s.handleBalance()
		span.SetOperationName("msg_handler: handle cmd `/balance`")
	default:
		resp = "не знаю эту команду"
	}
	return s.tgClient.SendMessage(resp, msg.UserID)
}

var errWrongFormat = errors.New("wrong format")

func handleF(span opentracing.Span, spanCtx context.Context, strs []string, count int, handler func(context.Context, []string) (string, error)) string {
	if count != len(strs) {
		return errWrongFormat.Error()
	}

	if r, err := handler(spanCtx, strs); err != nil {
		ext.Error.Set(span, true)
		return err.Error()
	} else {
		return r
	}
}
func genListMsg(els []string) string {
	var sb strings.Builder
	for i := 0; i < len(els); i++ {
		sb.WriteString(els[i])
		sb.WriteString("\n")
	}
	return sb.String()
}

func (s *MessageHandlerService) handleBalance() string {
	if v, err := s.stateService.GetBalance(); err != nil {
		return err.Error()
	} else {
		return fmt.Sprintf("%v rub", v)
	}
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
		els[i] = strconv.Itoa(allCats[i].Id) + " - " + allCats[i].Name
	}

	return genListMsg(els)
}

func (s *MessageHandlerService) handleAdd(ctx context.Context, tokens []string) (string, error) {
	catStr := tokens[1]
	sumStr := tokens[2]
	dtStr := tokens[3]
	var balanceAfter decimal.Decimal

	if cat, err := strconv.Atoi(catStr); err != nil {
		return "", errors.New("category must be a number")
	} else if sum, err := decimal.NewFromString(sumStr); err != nil {
		return "", errors.New("sum  must be a number")
	} else if dt, err := time.Parse(dtTemplate, dtStr); err != nil {
		return "", errors.New("wrong date format")
	} else if balanceAfter, err = s.spendingService.SaveTx(ctx, model.NewSpending(sum, cat, dt)); err != nil {
		return "", err
	}
	return fmt.Sprintf("added, current balance: %v", balanceAfter), nil
}

func parseReportReq(spanCtx context.Context, strs []string) (time.Time, time.Time, error) {
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
		return time.Time{}, time.Time{}, errWrongFormat
	}
	return startAt, endAt, nil
}
func (s *MessageHandlerService) handleReport(spanCtx context.Context, strs []string) (string, error) {
	startAt, endAt, err := parseReportReq(spanCtx, strs)
	if err != nil {
		return "", err
	}

	if data, c, err := s.spendingService.GetStatsBy(spanCtx, startAt, endAt); err != nil {
		return "", err
	} else {
		return formatStats(spanCtx, startAt, endAt, data, c), nil
	}
}

func (s *MessageHandlerService) handleReportAsync(spanCtx context.Context, userId int64, strs []string) (string, error) {
	startAt, endAt, err := parseReportReq(spanCtx, strs)
	if err != nil {
		return "", err
	}
	if err := s.reportProducer.Send(model.NewReportRequest(userId, startAt, endAt)); err != nil {
		return "", err
	}
	return "calculating report...", err
}

func (s *MessageHandlerService) reportResultListen(reportResultCh <-chan *model.Report) {
	for result := range reportResultCh {
		if err := s.tgClient.SendMessage(formatStats(context.Background(), result.Start, result.End, result.Data, ""), result.UserId); err != nil {
			Log.Error("failed to send report request", zap.Error(err))
		}
	}
}

func (s *MessageHandlerService) handleCurrencyChange(ctx context.Context, strs []string) (string, error) {
	if err := s.currencyService.UpdateCurrentCurrency(strs[1]); err != nil {
		return "", err
	}
	return "successfully changed", nil
}

func formatStats(spanCtx context.Context, start time.Time, end time.Time, r map[string]decimal.Decimal, currencyCode string) string {
	span, _ := opentracing.StartSpanFromContext(spanCtx, "msg_handler: formatStats response")
	defer span.Finish()

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
