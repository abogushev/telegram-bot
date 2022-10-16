package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Category int
type ReportType int

const (
	Other Category = iota
	Food
)

const (
	Week ReportType = iota
	Month
	Year
)

type Spending struct {
	Value    decimal.Decimal `db:"value"`
	Category Category        `db:"category"`
	Date     time.Time       `db:"date"`
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

func StrToCategory(str string) Category {
	switch str {
	case "food":
		return Food
	default:
		return Other
	}
}
