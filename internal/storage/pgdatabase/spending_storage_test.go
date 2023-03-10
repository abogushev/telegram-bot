package pgdatabase

import (
	"context"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/model"
)

func Test_Save(t *testing.T) {
	BeforeTest()
	storage := NewSpendingStorage(context.Background(), DB)
	tests := []struct {
		name     string
		data     model.Spending
		prepareF func()
		err      error
	}{
		{
			name: "save is ok",
			prepareF: func() {
			},
			err:  nil,
			data: model.Spending{Value: decimal.NewFromInt(1), CategoryId: 1, Date: time.Now()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			BeforeTest()
			tt.prepareF()
			err := storage.Save(tt.data)
			assert.ErrorIs(t, err, tt.err)
		})
	}
}

func Test_GetStatsBy(t *testing.T) {
	BeforeTest()
	storage := NewSpendingStorage(context.Background(), DB)

	tests := []struct {
		name     string
		startAt  time.Time
		endAt    time.Time
		data     model.ReportType
		prepareF func(start time.Time, end time.Time)
		checkF   func(map[string]decimal.Decimal, error)
	}{
		{
			name:    "report is ok",
			endAt:   time.Now(),
			startAt: time.Now().AddDate(0, 0, -7),
			prepareF: func(start time.Time, end time.Time) {
				DB.MustExec("insert into spendings(value, category_id, date) values(1, 1, $1)", start)
				DB.MustExec("insert into spendings(value, category_id, date) values(1, 1, $1)", start.AddDate(0, 0, 1))
				DB.MustExec("insert into spendings(value, category_id, date) values(1, 0, $1)", end)
			},
			data: model.Week,
			checkF: func(report map[string]decimal.Decimal, err error) {

				assert.True(t, report["food"].Equal(decimal.NewFromInt(1)))
				assert.True(t, report["other"].Equal(decimal.NewFromInt(2)))

				assert.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			BeforeTest()
			tt.prepareF(tt.startAt, tt.endAt)
			tt.checkF(storage.GetStatsBy(context.TODO(), tt.startAt, tt.endAt))
		})
	}
}
