package pgdatabase

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/model"
)

type dbStateStorage struct {
	ctx context.Context
	db  *sqlx.DB
}

func NewStateStorage(ctx context.Context, db *sqlx.DB) *dbStateStorage {
	return &dbStateStorage{ctx: ctx, db: db}
}

func(s *dbStateStorage) GetState() (model.State, error) {
	var state model.State
	if err := s.db.GetContext(s.ctx, &state, "select * from state"); err != nil {
		return model.State{}, err
	}
	return state, nil
}

func(s *dbStateStorage) UpdateBalance(v decimal.Decimal) error {
	if _, err := s.db.ExecContext(s.ctx, "update state set budget_balance = $1", v); err !=nil {
		return err
	}
	return nil
}

func(s *dbStateStorage) UpdateBalanceTx(tx *sqlx.Tx, v decimal.Decimal) error {
	_, err := tx.ExecContext(s.ctx, "update state set budget_balance = $1", v)
	return err
}

func(s *dbStateStorage) UpdateBalanceAndExpiresIn(v decimal.Decimal, t time.Time) error {
	if _, err := s.db.ExecContext(s.ctx, "update state set budget_balance = $1, budget_expires_in = $2", v, t); err !=nil {
		return err
	}
	return nil
}