package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type State struct {
	CurrentCurrencyCode string          `db:"current_currency_code"`
	BudgetValue         decimal.Decimal `db:"budget_value"`
	BudgetBalance       decimal.Decimal `db:"budget_balance"`
	BudgetExpiresIn     time.Time       `db:"budget_expires_in"`
}
