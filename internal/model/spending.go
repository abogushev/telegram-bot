package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Spending struct {
	Value      decimal.Decimal `db:"value"`
	CategoryId int             `db:"category_id"`
	Date       time.Time       `db:"date"`
}

func NewSpending(val decimal.Decimal, categoryId int, dt time.Time) Spending {
	return Spending{Value: val, CategoryId: categoryId, Date: dt}
}
