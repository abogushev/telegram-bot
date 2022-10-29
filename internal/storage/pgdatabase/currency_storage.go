package pgdatabase

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/model"
)

type dbCurrencyStorage struct {
	ctx context.Context
	db  *sqlx.DB
}

func NewCurrencyStorage(ctx context.Context, db *sqlx.DB) *dbCurrencyStorage {
	return &dbCurrencyStorage{ctx: ctx, db: db}
}

func (s *dbCurrencyStorage) GetCurrentCurrency(ctx context.Context) (model.Currency, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "currency_storage: GetCurrentCurrency")
	defer span.Finish()
	var c model.Currency
	q := "select code, ratio from currencies where code = (select current_currency_code from state)"
	if err := s.db.GetContext(s.ctx, &c, q); err != nil {
		return model.Currency{}, err
	}
	return c, nil
}

func (s *dbCurrencyStorage) GetCurrencies() ([]model.Currency, error) {
	cs := []model.Currency{}
	if err := s.db.SelectContext(s.ctx, &cs, "select code, ratio from currencies"); err != nil {
		return nil, err
	}
	return cs, nil
}

func (s *dbCurrencyStorage) UpdateCurrentCurrency(code string) error {
	_, err := s.db.ExecContext(s.ctx, "update state set current_currency_code = $1", code)
	return err
}

func (s *dbCurrencyStorage) UpdateCurrencies(newcrns []model.Currency) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}

	for i := 0; i < len(newcrns); i++ {
		if _, err := s.db.ExecContext(s.ctx, "insert into currencies values($1,$2) on conflict(code) do update set ratio = $2", newcrns[i].Code, newcrns[i].Ratio); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}

	return tx.Commit()
}
