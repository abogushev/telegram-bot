package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Category int
type ReportType int

const (
	Food Category = iota
	Other
)

const (
	Week ReportType = iota
	Month
	Year
)

type Spending struct {
	Value    decimal.Decimal
	Category Category
	Date     time.Time
}

func NewSpending(val decimal.Decimal, cat Category, dt time.Time) *Spending {
	return &Spending{Value: val, Category: cat, Date: dt}
}

func (c Category) String() string {
	r := ""
	switch c {
	case Food:
		r = "food"
	case Other:
		r = "other"
	}
	return r
}
