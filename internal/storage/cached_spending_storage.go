package storage

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/cache"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/model"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/time_util"
	"time"
)

type cacheI interface {
	Get(string) (string, error)
	Set(string, string, time.Duration) error
	Delete(string) error
}
type spendingStorageI interface {
	SaveTx(tx *sqlx.Tx, spending model.Spending) error
	GetStatsBy(context.Context, time.Time, time.Time) (map[string]decimal.Decimal, error)
}

var cacheKey = "report"

type CachedSpendingStorage struct {
	targetStorage spendingStorageI
	cache         cacheI
}

func NewCachedSpendingStorage(storage spendingStorageI, cacheI cacheI) *CachedSpendingStorage {
	return &CachedSpendingStorage{targetStorage: storage, cache: cacheI}
}

func (s *CachedSpendingStorage) SaveTx(tx *sqlx.Tx, spending model.Spending) error {
	if err := s.cache.Delete(cacheKey); err != nil {
		return err
	}
	return s.targetStorage.SaveTx(tx, spending)
}

func (s *CachedSpendingStorage) GetStatsBy(ctx context.Context, start time.Time, end time.Time) (map[string]decimal.Decimal, error) {
	cacheResult, err := s.cache.Get(cacheKey)
	if err != nil && err != cache.ErrNotFound {
		return nil, err
	}
	if err == cache.ErrNotFound {
		return s.refreshCache(ctx, start, end, make([]model.Report, 0))
	}

	reports, err := model.FromJSON(cacheResult)
	for i := 0; i < len(reports); i++ {
		if time_util.DatesEq(reports[i].Start, start) && time_util.DatesEq(reports[i].End, end) {
			return reports[i].Data, nil
		}
	}
	return s.refreshCache(ctx, start, end, reports)
}

func (s *CachedSpendingStorage) refreshCache(ctx context.Context, start, end time.Time, exists []model.Report) (map[string]decimal.Decimal, error) {
	result, err := s.targetStorage.GetStatsBy(ctx, start, end)
	if err != nil {
		return nil, err
	}
	newReport := model.NewReport(start, end, result)
	exists = append(exists, *newReport)

	js, err := model.ToJSON(exists)
	if err != nil {
		return nil, err
	}

	if err := s.cache.Set(cacheKey, js, time.Hour); err != nil {
		return nil, err
	}
	return result, nil
}
