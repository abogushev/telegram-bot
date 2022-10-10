package services

import (
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
	storage := mocks.NewMockSpendingService(ctrl)
	handlerService := NewMessageHandlerService(sender, storage)

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
	storage := mocks.NewMockSpendingService(ctrl)
	handlerService := NewMessageHandlerService(sender, storage)

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
	storage := mocks.NewMockSpendingService(ctrl)
	handlerService := NewMessageHandlerService(sender, storage)

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
	storage := mocks.NewMockSpendingService(ctrl)
	handlerService := NewMessageHandlerService(sender, storage)

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
	storage := mocks.NewMockSpendingService(ctrl)
	handlerService := NewMessageHandlerService(sender, storage)

	err := handlerService.HandleMsg(&model.Message{
		Text:   "/add 1 q 01-01-2000",
		UserID: 123,
	})

	assert.NoError(t, err)
}

func Test_OnAdd_shouldAnswerErrOnBadCatValue(t *testing.T) {
	ctrl := gomock.NewController(t)

	sender := mocks.NewMockMessageSender(ctrl)
	sender.EXPECT().SendMessage("wrong category", int64(123))
	storage := mocks.NewMockSpendingService(ctrl)
	handlerService := NewMessageHandlerService(sender, storage)

	err := handlerService.HandleMsg(&model.Message{
		Text:   "/add 99 1 01-01-2000",
		UserID: 123,
	})

	assert.NoError(t, err)
}

func Test_OnAdd_shouldAnswerErrOnBadDtFormat(t *testing.T) {
	ctrl := gomock.NewController(t)

	sender := mocks.NewMockMessageSender(ctrl)
	sender.EXPECT().SendMessage("wrong date format", int64(123))
	storage := mocks.NewMockSpendingService(ctrl)
	handlerService := NewMessageHandlerService(sender, storage)

	err := handlerService.HandleMsg(&model.Message{
		Text:   "/add 99 1 qwe",
		UserID: 123,
	})

	assert.NoError(t, err)
}

func Test_OnAdd_shouldSaveSuccessfull(t *testing.T) {
	ctrl := gomock.NewController(t)

	sender := mocks.NewMockMessageSender(ctrl)
	sender.EXPECT().SendMessage("added", int64(123))
	storage := mocks.NewMockSpendingService(ctrl)
	dt, _ := time.Parse("02-01-2006", "01-01-2000")
	storage.EXPECT().Save(model.NewSpending(decimal.NewFromInt(1), model.Category(1), dt))
	handlerService := NewMessageHandlerService(sender, storage)

	err := handlerService.HandleMsg(&model.Message{
		Text:   "/add 1 1 01-01-2000",
		UserID: 123,
	})

	assert.NoError(t, err)
}

func Test_OnAdd_shouldReportSuccessfull(t *testing.T) {
	ctrl := gomock.NewController(t)
	response := "from: 01-01-2000, to: 07-01-2000\nfood - 1 ₽\nother - 2 ₽\n"
	sender := mocks.NewMockMessageSender(ctrl)
	sender.EXPECT().SendMessage(response, int64(123))
	storage := mocks.NewMockSpendingService(ctrl)
	start, _ := time.Parse("02-01-2006", "01-01-2000")
	end, _ := time.Parse("02-01-2006", "07-01-2000")
	reportData := make(map[model.Category]decimal.Decimal)
	reportData[model.Food] = decimal.NewFromInt(1)
	reportData[model.Other] = decimal.NewFromInt(2)
	storage.EXPECT().GetStatsBy(model.Week).Return(start, end, reportData, model.RUB, nil)
	handlerService := NewMessageHandlerService(sender, storage)

	err := handlerService.HandleMsg(&model.Message{
		Text:   "/report w",
		UserID: 123,
	})

	assert.NoError(t, err)
}
