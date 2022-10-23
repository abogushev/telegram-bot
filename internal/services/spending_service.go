package services

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/model"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/storage/pgdatabase"
)

type spendingStorageI interface {
	Save(model.Spending) error
	SaveTx(tx *sqlx.Tx, spending model.Spending) error
	GetStatsBy(time.Time, time.Time) (map[string]decimal.Decimal, error)
}

type currencyServiceI interface {
	GetCurrentCurrency() (model.Currency, error)
}
type stateServiceI interface {
	DecreaseBalanceTx(tx *sqlx.Tx, v decimal.Decimal) (decimal.Decimal, error)
}
type spendingService struct {
	spendingStorage spendingStorageI
	currencyService currencyServiceI
	stateService    stateServiceI
}

func NewSpendingService(spendingStorage spendingStorageI, currencyService currencyServiceI, stateServiceTx stateServiceI) *spendingService {
	return &spendingService{spendingStorage, currencyService, stateServiceTx}
}

func (s *spendingService) Save(spending model.Spending) error {
	if cur, err := s.currencyService.GetCurrentCurrency(); err != nil {
		return err
	} else {
		spending.Value = spending.Value.Div(cur.Ratio)
		return s.spendingStorage.Save(spending)
	}
}

func (s *spendingService) SaveTx(spending model.Spending) (decimal.Decimal, error) {
	if cur, err := s.currencyService.GetCurrentCurrency(); err != nil {
		return decimal.Decimal{}, err
	} else {
		spending.Value = spending.Value.Div(cur.Ratio)
	}

	var balanceAfter decimal.Decimal
	err := pgdatabase.RunInTx(
		func(tx *sqlx.Tx) error {
			var err error
			log.Println("start DecreaseBalanceTx")
			balanceAfter, err = s.stateService.DecreaseBalanceTx(tx, spending.Value)
			log.Println("end DecreaseBalanceTx")
			return err
		},
		func(tx *sqlx.Tx) error {
			log.Println("start SaveTx")
			err := s.spendingStorage.SaveTx(tx, spending)
			log.Println("end SaveTx")
			return err
		},
	)
	return balanceAfter, err
}

func (s *spendingService) GetStatsBy(start, end time.Time) (map[string]decimal.Decimal, string, error) {
	data, err := s.spendingStorage.GetStatsBy(start, end)
	if err != nil {
		return nil, "", err
	}
	rs := make(map[string]decimal.Decimal)
	ct, err := s.currencyService.GetCurrentCurrency()
	if err != nil {
		return nil, "", err
	}
	for k, v := range data {
		rs[k] = ct.Ratio.Mul(v)
	}
	return rs, ct.Code, nil
}
