package services

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	mocks "gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/mocks/services"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/model"
)

func Test_OnStartCommand_ShouldAnswerWithIntroMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	sender := mocks.NewMockMessageSender(ctrl)

	handlerService := NewMessageHandlerService(
		sender,
		mocks.NewMockSpendingService(ctrl),
		mocks.NewMockCurrencyService(ctrl),
		mocks.NewMockCategoryService(ctrl),
		mocks.NewMockStateService(ctrl),
	)

	sender.EXPECT().SendMessage("hello", int64(123))

	err := handlerService.HandleMsg(&model.Message{
		Text:   "/start",
		UserID: 123,
	})

	assert.NoError(t, err)
}

func Test_OnUnknownCommand_ShouldAnswerWithHelpMessage(t *testing.T) {
	ctrl := gomock.NewController(t)

	sender := mocks.NewMockMessageSender(ctrl)
	sender.EXPECT().SendMessage("не знаю эту команду", int64(123))
	handlerService := NewMessageHandlerService(
		sender,
		mocks.NewMockSpendingService(ctrl),
		mocks.NewMockCurrencyService(ctrl),
		mocks.NewMockCategoryService(ctrl),
		mocks.NewMockStateService(ctrl),
	)

	err := handlerService.HandleMsg(&model.Message{
		Text:   "some text",
		UserID: 123,
	})

	assert.NoError(t, err)
}

func Test_OnAdd_shouldAnswerErrOnWrongCountOfTokens(t *testing.T) {
	ctrl := gomock.NewController(t)

	sender := mocks.NewMockMessageSender(ctrl)
	sender.EXPECT().SendMessage("wrong format", int64(123))
	handlerService := NewMessageHandlerService(
		sender,
		mocks.NewMockSpendingService(ctrl),
		mocks.NewMockCurrencyService(ctrl),
		mocks.NewMockCategoryService(ctrl),
		mocks.NewMockStateService(ctrl),
	)

	err := handlerService.HandleMsg(&model.Message{
		Text:   "/add 1 01-01-2000",
		UserID: 123,
	})

	assert.NoError(t, err)
}

func Test_OnAdd_shouldAnswerErrOnNonNumberCatValue(t *testing.T) {
	ctrl := gomock.NewController(t)

	sender := mocks.NewMockMessageSender(ctrl)
	sender.EXPECT().SendMessage("category must be a number", int64(123))

	handlerService := NewMessageHandlerService(
		sender,
		mocks.NewMockSpendingService(ctrl),
		mocks.NewMockCurrencyService(ctrl),
		mocks.NewMockCategoryService(ctrl),
		mocks.NewMockStateService(ctrl),
	)

	err := handlerService.HandleMsg(&model.Message{
		Text:   "/add q 1 01-01-2000",
		UserID: 123,
	})

	assert.NoError(t, err)
}

func Test_OnAdd_shouldAnswerErrOnNonNumberSumValue(t *testing.T) {
	ctrl := gomock.NewController(t)

	sender := mocks.NewMockMessageSender(ctrl)
	sender.EXPECT().SendMessage("sum  must be a number", int64(123))

	handlerService := NewMessageHandlerService(
		sender,
		mocks.NewMockSpendingService(ctrl),
		mocks.NewMockCurrencyService(ctrl),
		mocks.NewMockCategoryService(ctrl),
		mocks.NewMockStateService(ctrl),
	)

	err := handlerService.HandleMsg(&model.Message{
		Text:   "/add 1 q 01-01-2000",
		UserID: 123,
	})

	assert.NoError(t, err)
}

func Test_OnAdd_shouldAnswerErrOnBadDtFormat(t *testing.T) {
	ctrl := gomock.NewController(t)

	sender := mocks.NewMockMessageSender(ctrl)
	sender.EXPECT().SendMessage("wrong date format", int64(123))

	handlerService := NewMessageHandlerService(
		sender,
		mocks.NewMockSpendingService(ctrl),
		mocks.NewMockCurrencyService(ctrl),
		mocks.NewMockCategoryService(ctrl),
		mocks.NewMockStateService(ctrl),
	)

	err := handlerService.HandleMsg(&model.Message{
		Text:   "/add 99 1 qwe",
		UserID: 123,
	})

	assert.NoError(t, err)
}

func Test_OnAdd_shouldSaveSuccessfull(t *testing.T) {
	ctrl := gomock.NewController(t)

	sender := mocks.NewMockMessageSender(ctrl)
	sender.EXPECT().SendMessage("added, current balance: 0", int64(123))
	storage := mocks.NewMockSpendingService(ctrl)
	dt, _ := time.Parse("02-01-2006", "01-01-2000")
	storage.EXPECT().SaveTx(model.NewSpending(decimal.NewFromInt(1), 1, dt))
	handlerService := NewMessageHandlerService(
		sender,
		storage,
		mocks.NewMockCurrencyService(ctrl),
		mocks.NewMockCategoryService(ctrl),
		mocks.NewMockStateService(ctrl),
	)

	err := handlerService.HandleMsg(&model.Message{
		Text:   "/add 1 1 01-01-2000",
		UserID: 123,
	})

	assert.NoError(t, err)
}

func Test_OnAdd_shouldReportSuccessfull(t *testing.T) {
	ctrl := gomock.NewController(t)
	end := time.Now().Truncate(24 * time.Hour)
	start := end.AddDate(0, 0, -7)
	response := fmt.Sprintf("from: %v, to: %v\nfood - 1 rub\nother - 2 rub\n", start.Format("02-01-2006"), end.Format("02-01-2006"))
	sender := mocks.NewMockMessageSender(ctrl)
	sender.EXPECT().SendMessage(response, int64(123))
	storage := mocks.NewMockSpendingService(ctrl)

	reportData := make(map[string]decimal.Decimal)
	reportData["food"] = decimal.NewFromInt(1)
	reportData["other"] = decimal.NewFromInt(2)
	storage.EXPECT().GetStatsBy(start, end).Return(reportData, "rub", nil)
	handlerService := NewMessageHandlerService(
		sender,
		storage,
		mocks.NewMockCurrencyService(ctrl),
		mocks.NewMockCategoryService(ctrl),
		mocks.NewMockStateService(ctrl),
	)

	err := handlerService.HandleMsg(&model.Message{
		Text:   "/report w",
		UserID: 123,
	})

	assert.NoError(t, err)
}
