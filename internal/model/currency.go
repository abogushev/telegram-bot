package model

import (
	"errors"

	"github.com/shopspring/decimal"
)

var ErrWrongCurrency = errors.New("wrong currency type")

type Currency struct {
	Code  string          `db:"code"`
	Ratio decimal.Decimal `db:"ratio"`
}

func NewCurrency(code string, ratio decimal.Decimal) *Currency {
	return &Currency{code, ratio}
}
