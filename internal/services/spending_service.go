package services

import (
	"time"

	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/model"
)

type spendingStorage interface {
	Save(*model.Spending) error
	GetStatsBy(time.Time, time.Time) (map[model.Category]decimal.Decimal, error)
}

type currencyStorage interface {
	UpdateCurrentType(model.CurrencyType) error
	GetCurrencyType() (model.CurrencyType, error)
	GetCurrencyRatioToRUB(c model.CurrencyType) (decimal.Decimal, error)
}

type spendingService struct {
	spendingStorage spendingStorage
	currencyStorage currencyStorage
}

func NewSpendingService(spendingStorage spendingStorage, currencyStorage currencyStorage) *spendingService {
	return &spendingService{spendingStorage, currencyStorage}
}

func (s *spendingService) UpdateCurrentType(t model.CurrencyType) error {
	return s.currencyStorage.UpdateCurrentType(t)
}

func (s *spendingService) Save(spending *model.Spending) error {
	if cnvV, err := s.ConvertFromCurrentCurrencyToRUB(spending.Value); err != nil {
		return err
	} else {
		spending.Value = cnvV
		return s.spendingStorage.Save(spending)
	}
}

func (s *spendingService) cnvrtF(value decimal.Decimal, cnv func(ratio, v decimal.Decimal) decimal.Decimal) (decimal.Decimal, error) {
	ct, err := s.currencyStorage.GetCurrencyType()
	if err != nil {
		return decimal.Decimal{}, err
	}
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

func (s *spendingService) GetStatsBy(start, end time.Time) (map[model.Category]decimal.Decimal, model.CurrencyType, error) {
	data, err := s.spendingStorage.GetStatsBy(start, end)
	if err != nil {
		return nil, model.Undefined, err
	}
	ct, err := s.currencyStorage.GetCurrencyType()
	if err != nil {
		return nil, model.Undefined, err
	}
	if ct == model.RUB {
		return data, ct, nil
	}
	if ratio, err := s.currencyStorage.GetCurrencyRatioToRUB(ct); err != nil {
		return nil, ct, err
	} else {
		rs := make(map[model.Category]decimal.Decimal)
		for k, v := range data {
			rs[k] = ratio.Mul(v)
		}
		return rs, ct, nil
	}
}
