package model

import "errors"

type CurrencyType int

const (
	USD CurrencyType = iota
	CNY
	EUR
	RUB
	Undefined
)

var ErrWrongCurrencyType = errors.New("wrong currency type")

func ParseCurrencyType(t int) (CurrencyType, error) {
	switch CurrencyType(t) {
	case USD | CNY | EUR | RUB:
		return CurrencyType(t), nil
	default:
		return Undefined, ErrWrongCurrencyType
	}
}
