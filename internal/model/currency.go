package model

import "errors"

type CurrencyType int

const (
	Undefined CurrencyType = iota
	USD
	CNY
	EUR
	RUB
)

func (c CurrencyType) String() string {
	switch c {
	case USD:
		return "$"
	case CNY:
		return "¥"
	case EUR:
		return "€"
	case RUB:
		return "₽"
	default:
		return "?"
	}
}

var AllNonBaseCurrenciesType = []CurrencyType{USD, CNY, EUR}
var ErrWrongCurrencyType = errors.New("wrong currency type")

func ParseCurrencyType(t int) (CurrencyType, error) {
	switch {
	case t >= 0 && t < int(RUB):
		return CurrencyType(t), nil
	default:
		return Undefined, ErrWrongCurrencyType
	}
}
