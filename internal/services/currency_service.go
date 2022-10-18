package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/model"
)

var url = "https://api.currencyapi.com/v3/latest?apikey=dO62Nn8Y3f18mpvbN6ypoaBrzEtKF8Fkd8bdavYy&currencies=EUR%2CUSD%2CCNY&base_currency=RUB"

type currencyService struct {
	jobMutex          sync.Once
	currencies        map[string]model.Currency
	currenciesM       sync.RWMutex
	currenciesStorage currenciesStorage
	currentCurrency   model.Currency
	currentCurrencyM  sync.RWMutex
}

func NewCurrencyService(currenciesStorage currenciesStorage) (*currencyService, error) {
	currentCurrency, err := currenciesStorage.GetCurrentCurrency()
	if err != nil {
		return nil, err
	}
	currencies, err := currenciesStorage.GetCurrencies()
	if err != nil {
		return nil, err
	}

	if len(currencies) > 1 {
		mcurrencies := make(map[string]model.Currency)
		for i := 0; i < len(currencies); i++ {
			mcurrencies[currencies[i].Code] = currencies[i]
		}
		log.Println("CURRENCIES:", mcurrencies)
		return &currencyService{
			currencies:        mcurrencies,
			currentCurrency:   currentCurrency,
			currenciesStorage: currenciesStorage,
		}, nil
	} else {
		cs := &currencyService{
			currencies:        map[string]model.Currency{},
			currentCurrency:   currentCurrency,
			currenciesStorage: currenciesStorage,
		}
		if err := cs.updateCurrencies(context.Background()); err != nil {
			return nil, err
		}
		log.Println("CURRENCIES AFTER LOAD:", cs.currencies)
		return cs, nil
	}
}

type currenciesStorage interface {
	GetCurrentCurrency() (model.Currency, error)
	GetCurrencies() ([]model.Currency, error)
	UpdateCurrencies([]model.Currency) error
	UpdateCurrentCurrency(name string) error
}

func (cs *currencyService) GetAll() []model.Currency {
	cs.currenciesM.RLock()
	defer cs.currenciesM.RUnlock()
	result := make([]model.Currency, 0, len(cs.currencies))
	for _, v := range cs.currencies {
		result = append(result, v)
	}
	log.Println("RETURN CURRENCIES:", result)
	return result
}

func (cs *currencyService) CheckCurrencyCode(code string) bool {
	cs.currenciesM.RLock()
	defer cs.currenciesM.RUnlock()
	_, ok := cs.currencies[code]
	return ok
}

func (cs *currencyService) GetCurrentCurrency() model.Currency {
	cs.currentCurrencyM.RLock()
	defer cs.currentCurrencyM.RUnlock()
	return cs.currentCurrency
}

func (cs *currencyService) UpdateCurrentCurrency(newCur string) error {
	cs.currenciesM.RLock()
	currency, ok := cs.currencies[newCur]
	cs.currenciesM.RUnlock()
	if !ok {
		return model.ErrWrongCurrency
	}

	cs.currentCurrencyM.RLock()
	defer cs.currentCurrencyM.RUnlock()

	if currency.Code == cs.currentCurrency.Code {
		return nil
	}

	if err := cs.currenciesStorage.UpdateCurrentCurrency(currency.Code); err != nil {
		return err
	}
	cs.currentCurrency = currency
	return nil
}

func (s *currencyService) RunUpdateCurrenciesDaemon(ctx context.Context, updateInterval time.Duration) {
	go s.jobMutex.Do(func() {
		ticker := time.NewTicker(updateInterval)

		for {
			select {
			case <-ticker.C:
				if err := s.updateCurrencies(ctx); err != nil {
					fmt.Printf("error on update currencies, %v\n", err)
				}
			case <-ctx.Done():
				fmt.Printf("cancel update currencies job")
				return
			}
		}
	})
}

func (s *currencyService) updateCurrencies(ctx context.Context) error {
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
	if len(rs.Data) == 0 {
		return nil
	}
	arr := make([]model.Currency, 0, len(rs.Data))
	mapCt := make(map[string]model.Currency)
	for code, v := range rs.Data {
		m := *model.NewCurrency(code, v.Value)
		arr = append(arr, m)
		mapCt[m.Code] = m
	}

	s.currenciesM.Lock()
	defer s.currenciesM.Unlock()
	if err := s.currenciesStorage.UpdateCurrencies(arr); err != nil {
		return err
	}
	s.currencies = mapCt
	fmt.Println("CURRENCIES UPDATED: ", s.currencies)
	return nil
}

type dataResponse struct {
	Data map[string]currencyResponse `json:"data"`
}
type currencyResponse struct {
	Value decimal.Decimal `json:"value"`
}
