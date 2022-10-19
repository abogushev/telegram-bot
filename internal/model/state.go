package model

import (
	"github.com/shopspring/decimal"
	"time"
)

type State struct {
	CurrentCurrencyCode string `db:"current_currency_code"`
	BudgetValue decimal.Decimal `db:"budget_value"`
	BudgetBalance decimal.Decimal `db:"budget_balance"`
	BudgetExpiresIn time.Time `db:"budget_expires_in"`
}