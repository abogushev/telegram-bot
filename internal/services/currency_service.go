package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/model"
)

var url = "https://api.currencyapi.com/v3/latest?apikey=dO62Nn8Y3f18mpvbN6ypoaBrzEtKF8Fkd8bdavYy&currencies=EUR%2CUSD%2CCNY&base_currency=RUB"

type CurrencyService struct {
	jobMutex sync.Once
}

func NewCurrencyService() *CurrencyService {
	return &CurrencyService{}
}

type currenciesStorage interface {
	UpdateCurrencies(map[model.CurrencyType]decimal.Decimal) error
}

func (s *CurrencyService) RunUpdateCurrenciesDaemon(ctx context.Context, updateInterval time.Duration, currenciesStorage currenciesStorage) {
	go s.jobMutex.Do(func() {
		ticker := time.NewTicker(updateInterval)

		for {
			select {
			case <-ticker.C:
				if err := updateCurrencies(ctx, currenciesStorage); err != nil {
					fmt.Printf("error on update currencies, %v\n", err)
				}
			case <-ctx.Done():
				fmt.Printf("cancel update currencies job")
				return
			}
		}
	})
}

func updateCurrencies(ctx context.Context, currenciesStorage currenciesStorage) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return err
	}
	req.Header.Add("accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	rs := &dataResponse{}
	if err := json.Unmarshal(body, &rs); err != nil {
		return err
	}
	currencies := make(map[model.CurrencyType]decimal.Decimal)

	for ctype, cur := range rs.Data {
		switch ctype {
		case "USD":
			currencies[model.USD] = cur.Value
		case "CNY":
			currencies[model.CNY] = cur.Value
		case "EUR":
			currencies[model.EUR] = cur.Value
		}
	}

	return currenciesStorage.UpdateCurrencies(currencies)
}

type dataResponse struct {
	Data map[string]currencyResponse `json:"data"`
}
type currencyResponse struct {
	Value decimal.Decimal `json:"value"`
}
