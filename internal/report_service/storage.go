package report_service

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"time"
)

type ReportStorage struct {
	DB *sqlx.DB
}

func (s *ReportStorage) getStatsBy(ctx context.Context, startAt, endAt time.Time) (map[string]float64, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "storage: getting report")
	defer span.Finish()

	results := []struct {
		Name  string  `db:"name"`
		Value float64 `db:"value"`
	}{}

	q := "select categories.name as name, sum(spendings.value) as value from spendings inner join categories on spendings.category_id = categories.id where date between $1 and $2 group by categories.name"
	if err := s.DB.Select(&results, q, startAt, endAt); err != nil {
		ext.Error.Set(span, true)
		return nil, err
	}

	r := make(map[string]float64)
	for i := 0; i < len(results); i++ {
		r[results[i].Name] = results[i].Value
	}

	return r, nil
}
