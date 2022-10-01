package storage

import (
	"time"

	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/model"
)

type InMemoryStorage struct {
	data map[model.Category][]*model.Spending
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{make(map[model.Category][]*model.Spending)}
}

func (s *InMemoryStorage) Save(spending *model.Spending) {
	s.data[spending.Category] = append(s.data[spending.Category], spending)
}

func (s *InMemoryStorage) GetStatsBy(rt model.ReportType) (time.Time, time.Time, map[model.Category]decimal.Decimal) {
	endAt := time.Now()
	var startAt time.Time
	switch rt {
	case model.Week:
		startAt = endAt.AddDate(0, 0, -7)
	case model.Month:
		startAt = endAt.AddDate(0, -1, 0)
	case model.Year:
		startAt = endAt.AddDate(-1, 0, 0)
	}
	return startAt, endAt, s.groupBy(startAt, endAt)
}

func (s *InMemoryStorage) groupBy(startAt time.Time, endAt time.Time) map[model.Category]decimal.Decimal {
	result := make(map[model.Category]decimal.Decimal)
	for cat, ms := range s.data {
		for i := 0; i < len(ms); i++ {
			if (ms[i].Date.After(startAt) && ms[i].Date.Before(endAt)) ||
				(ms[i].Date.Equal(startAt) || ms[i].Date.Equal(endAt)) {
				result[cat] = result[cat].Add(ms[i].Value)
			}
		}
	}
	return result
}
