package services

import (
	"time"

	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/model"
)

type spendingStorageI interface {
	Save(model.Spending) error
	GetStatsBy(time.Time, time.Time) (map[string]decimal.Decimal, error)
}

type currencyServiceI interface {
	GetCurrentCurrency() model.Currency
}

type spendingService struct {
	spendingStorage spendingStorageI
	currencyService currencyServiceI
}

func NewSpendingService(spendingStorage spendingStorageI, currencyService currencyServiceI) *spendingService {
	return &spendingService{spendingStorage, currencyService}
}

func (s *spendingService) Save(spending model.Spending) error {
	spending.Value = spending.Value.Div(s.currencyService.GetCurrentCurrency().Ratio)
	return s.spendingStorage.Save(spending)
}

func (s *spendingService) GetStatsBy(start, end time.Time) (map[string]decimal.Decimal, string, error) {
	data, err := s.spendingStorage.GetStatsBy(start, end)
	if err != nil {
		return nil, "", err
	}
	rs := make(map[string]decimal.Decimal)
	ct := s.currencyService.GetCurrentCurrency()
	for k, v := range data {
		rs[k] = ct.Ratio.Mul(v)
	}
	return rs, ct.Code, nil
}
