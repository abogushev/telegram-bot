package services

import (
	"time"

	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/model"
)

type spendingStorage interface {
	Save(*model.Spending)
	GetStatsBy(model.ReportType) (time.Time, time.Time, map[model.Category]decimal.Decimal)
}

type currencyStorage interface {
	UpdateCurrentType(model.CurrencyType)
	GetCurrencyType() model.CurrencyType
	GetCurrencyRatioToRUB(c model.CurrencyType) (decimal.Decimal, error)
}

type spendingService struct {
	spendingStorage spendingStorage
	currencyStorage currencyStorage
}

func NewSpendingService(spendingStorage spendingStorage, currencyStorage currencyStorage) *spendingService {
	return &spendingService{spendingStorage, currencyStorage}
}

func (s *spendingService) UpdateCurrentType(t model.CurrencyType) {
	s.currencyStorage.UpdateCurrentType(t)
}

func (s *spendingService) Save(spending *model.Spending) error {
	if cnvV, err := s.ConvertFromCurrentCurrencyToRUB(spending.Value); err != nil {
		return err
	} else {
		spending.Value = cnvV
		s.spendingStorage.Save(spending)
		return nil
	}
}

func (s *spendingService) cnvrtF(value decimal.Decimal, cnv func(ratio, v decimal.Decimal) decimal.Decimal) (decimal.Decimal, error) {
	ct := s.currencyStorage.GetCurrencyType()
	if ct == model.RUB {
		return value, nil
	}
	if ratio, err := s.currencyStorage.GetCurrencyRatioToRUB(ct); err != nil {
		return decimal.Decimal{}, err
	} else {
		return cnv(ratio, value), nil
	}
}

func (s *spendingService) ConvertFromCurrentCurrencyToRUB(value decimal.Decimal) (decimal.Decimal, error) {
	return s.cnvrtF(value, func(ratio, v decimal.Decimal) decimal.Decimal { return value.Div(ratio) })
}

func (s *spendingService) ConvertRUBToCurrentCurrencyType(value decimal.Decimal) (decimal.Decimal, error) {
	return s.cnvrtF(value, func(ratio, v decimal.Decimal) decimal.Decimal { return ratio.Mul(value) })
}

func (s *spendingService) GetStatsBy(t model.ReportType) (time.Time, time.Time, map[model.Category]decimal.Decimal, model.CurrencyType, error) {
	start, end, data := s.spendingStorage.GetStatsBy(t)
	ct := s.currencyStorage.GetCurrencyType()
	if ct == model.RUB {
		return start, end, data, ct, nil
	}
	if ratio, err := s.currencyStorage.GetCurrencyRatioToRUB(ct); err != nil {
		return time.Time{}, time.Time{}, nil, ct, err
	} else {
		rs := make(map[model.Category]decimal.Decimal)
		for k, v := range data {
			rs[k] = ratio.Mul(v)
		}
		return start, end, rs, ct, nil
	}
}
