package services

import (
	"context"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/model"
)

type stateStorage interface {
	GetState() (model.State, error)
	DecreaseBalanceTx(*sqlx.Tx, decimal.Decimal) (decimal.Decimal, error)
	UpdateBalanceAndExpiresIn(time.Time) error
}
type stateService struct {
	stateStorage stateStorage
}

func NewStateService(storage stateStorage, ctx context.Context) (*stateService, error) {
	state, err := storage.GetState()
	if err != nil {
		return nil, err
	}
	service := &stateService{stateStorage: storage}

	go service.runJob(ctx, time.Until(state.BudgetExpiresIn))

	return service, nil
}

func (s *stateService) GetBalance() (decimal.Decimal, error) {
	state, err := s.stateStorage.GetState()
	if err != nil {
		return decimal.Decimal{}, err
	}
	return state.BudgetBalance, nil
}

func (s *stateService) DecreaseBalanceTx(tx *sqlx.Tx, v decimal.Decimal) (decimal.Decimal, error) {
	log.Println("start UpdateBalanceTx")
	result, err := s.stateStorage.DecreaseBalanceTx(tx, v)
	if err != nil {
		return decimal.Decimal{}, err
	}
	log.Println("end UpdateBalanceTx")
	return result, nil
}

func (s *stateService) runJob(ctx context.Context, nextTriggerTime time.Duration) {
	timer := time.NewTimer(nextTriggerTime)

	for {
		select {
		case <-timer.C:
			nextTime := time.Now().AddDate(0, 1, 0).Truncate(24 * time.Hour)

			timer = time.NewTimer(time.Until(nextTime))

			if err := s.stateStorage.UpdateBalanceAndExpiresIn(nextTime); err != nil {
				log.Printf("error on update state, %v\n", err)
			}
		case <-ctx.Done():
			log.Printf("cancel update state job")
			return
		}
	}
}
