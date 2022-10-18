package pgdatabase

import (
	"context"

	"github.com/jmoiron/sqlx"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/model"
)

type dbCategoryStorage struct {
	ctx context.Context
	db  *sqlx.DB
}

func NewCategoryStorage(ctx context.Context, db *sqlx.DB) *dbCategoryStorage {
	return &dbCategoryStorage{ctx: ctx, db: db}
}

func (s *dbCategoryStorage) GetAll() ([]model.Category, error) {
	r := []model.Category{}
	err := s.db.SelectContext(s.ctx, &r, "select id, name from categories")
	return r, err
}
