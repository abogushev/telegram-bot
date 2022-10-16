package pgdatabase

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/model"
)

type dbCurrencyStorage struct {
	ctx context.Context
	db  *sqlx.DB
}

func NewCurrencyStorage(ctx context.Context, db *sqlx.DB) *dbCurrencyStorage {
	return &dbCurrencyStorage{ctx: ctx, db: db}
}

func (s *dbCurrencyStorage) UpdateCurrentType(ctype model.CurrencyType) error {
	code, err := model.CurrencyTypeToStr(ctype)
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(s.ctx, "update state set current_currency_code = $1", code)
	return err
}

func (s *dbCurrencyStorage) GetCurrencyType() (model.CurrencyType, error) {
	var c string
	if err := s.db.GetContext(s.ctx, &c, "select current_currency_code from state"); err != nil {
		return model.Undefined, err
	}
	return model.StrToCurrencyType(c)
}

func (s *dbCurrencyStorage) GetCurrencyRatioToRUB(c model.CurrencyType) (decimal.Decimal, error) {
	var r decimal.Decimal
	code, err := model.CurrencyTypeToStr(c)
	if err != nil {
		return decimal.Decimal{}, err
	}
	err = s.db.GetContext(s.ctx, &r, "select ratio from currencies where code = $1", code)
	if err != nil {
		return decimal.Decimal{}, err
	}
	return r, nil
}

func (s *dbCurrencyStorage) UpdateCurrencies(newcrns map[model.CurrencyType]decimal.Decimal) error {
	upds := make(map[string]decimal.Decimal, 0)
	for k, v := range newcrns {
		s, err := model.CurrencyTypeToStr(k)
		if err != nil {
			return errors.Wrapf(err, "wrong сurrency type: %v", int(k))
		}
		upds[s] = v
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}

	for k, v := range upds {
		if _, err := s.db.ExecContext(s.ctx, "insert into currencies values($1,$2) on conflict(code) do update set ratio = $2", k, v); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}

	return tx.Commit()
}
