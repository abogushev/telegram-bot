package model

import (
	"errors"
	"strings"
)

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

func StrToCurrencyType(s string) (CurrencyType, error) {
	switch strings.ToLower(s) {
	case "usd":
		return USD, nil
	case "eur":
		return EUR, nil
	case "cny":
		return CNY, nil
	case "rub":
		return RUB, nil
	default:
		return Undefined, ErrWrongCurrencyType
	}
}

func CurrencyTypeToStr(c CurrencyType) (string, error) {
	switch c {
	case USD:
		return "usd", nil
	case CNY:
		return "cny", nil
	case EUR:
		return "eur", nil
	case RUB:
		return "rub", nil
	default:
		return "", ErrWrongCurrencyType
	}
}
