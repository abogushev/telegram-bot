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

func (s *dbSpendingStorage) Save(spending model.Spending) error {
	_, err := s.db.ExecContext(s.ctx, "insert into spendings(value, category_id, date) values($1,$2,$3)", spending.Value, spending.CategoryId, spending.Date)
	return err
}
func (s *dbSpendingStorage) SaveTx(tx *sqlx.Tx, spending model.Spending) error {
	if _, err := tx.ExecContext(s.ctx, "insert into spendings(value, category_id, date) values($1,$2,$3)", spending.Value, spending.CategoryId, spending.Date); err != nil {
		return err
	}
	return nil
}

func (s *dbSpendingStorage) GetStatsBy(startAt, endAt time.Time) (map[string]decimal.Decimal, error) {
	results := []struct {
		Name  string          `db:"name"`
		Value decimal.Decimal `db:"value"`
	}{}

	q := "select categories.name as name, sum(spendings.value) as value from spendings inner join categories on spendings.category_id = categories.id where date between $1 and $2 group by categories.name"
	if err := s.db.Select(&results, q, startAt, endAt); err != nil {
		return nil, err
	}

	r := make(map[string]decimal.Decimal)
	for i := 0; i < len(results); i++ {
		r[results[i].Name] = results[i].Value
	}

	return r, nil
}
