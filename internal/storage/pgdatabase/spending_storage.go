package pgdatabase

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/model"
)

type dbSpendingStorage struct {
	ctx context.Context
	db  *sqlx.DB
}

func NewSpendingStorage(ctx context.Context, db *sqlx.DB) *dbSpendingStorage {
	return &dbSpendingStorage{ctx: ctx, db: db}
}

func (s *dbSpendingStorage) Save(spending *model.Spending) error {
	_, err := s.db.ExecContext(s.ctx, "insert into spendings(value, category, date) values($1,$2,$3)", spending.Value, spending.Category.String(), spending.Date)
	return err
}

func (s *dbSpendingStorage) GetStatsBy(startAt, endAt time.Time) (map[model.Category]decimal.Decimal, error) {
	results := []struct {
		Category string
		Value    decimal.Decimal
	}{}
	if err := s.db.Select(&results, "select category, sum(value) as value from spendings where date between $1 and $2 group by category;", startAt, endAt); err != nil {
		return map[model.Category]decimal.Decimal{}, err
	}

	r := make(map[model.Category]decimal.Decimal)
	for i := 0; i < len(results); i++ {
		r[model.StrToCategory(results[i].Category)] = results[i].Value
	}

	return r, nil
}
