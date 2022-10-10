package storage

import (
	"errors"
	"sync"

	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/model"
)

type currencyStorage struct {
	currencyType      model.CurrencyType
	currencyTypeMutex sync.RWMutex
	currencies        map[model.CurrencyType]decimal.Decimal
	currenciesMutex   sync.RWMutex
}

func NewCurrencyStorage() *currencyStorage {
	return &currencyStorage{
		currencyType:      model.RUB,
		currencyTypeMutex: sync.RWMutex{},
		currencies:        make(map[model.CurrencyType]decimal.Decimal),
		currenciesMutex:   sync.RWMutex{},
	}
}

func (m *currencyStorage) UpdateCurrentType(newType model.CurrencyType) {
	m.currencyTypeMutex.Lock()
	m.currencyType = newType
	m.currencyTypeMutex.Unlock()
}

var ErrCurrenciesNotInitialized = errors.New("currencies not initialized")

func (m *currencyStorage) GetCurrencyType() model.CurrencyType {
	m.currencyTypeMutex.RLock()
	defer m.currencyTypeMutex.RUnlock()
	return m.currencyType
}

func (m *currencyStorage) GetCurrencyRatioToRUB(c model.CurrencyType) (decimal.Decimal, error) {
	m.currenciesMutex.RLock()
	defer m.currenciesMutex.RUnlock()

	if v, ok := m.currencies[c]; !ok {
		return decimal.Decimal{}, ErrCurrenciesNotInitialized
	} else {
		return v, nil
	}
}

func (m *currencyStorage) UpdateCurrencies(newcrns map[model.CurrencyType]decimal.Decimal) {
	m.currenciesMutex.Lock()
	m.currencies = newcrns
	m.currenciesMutex.Unlock()
}
