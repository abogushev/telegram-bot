package services

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/model"
)

type stateStorage interface {
	GetState() (model.State, error)
	UpdateBalance(decimal.Decimal) error
	UpdateBalanceTx(*sqlx.Tx, decimal.Decimal) error
	UpdateBalanceAndExpiresIn(decimal.Decimal, time.Time) error
}
type stateService struct {
	stateStorage stateStorage
	state        model.State
	stateM       sync.RWMutex
}

func NewStateService(storage stateStorage, ctx context.Context) (*stateService, error) {
	state, err := storage.GetState()
	if err != nil {
		return nil, err
	}
	service := &stateService{stateStorage: storage, state: state, stateM: sync.RWMutex{}}

	go service.runJob(ctx, time.Now().Sub(state.BudgetExpiresIn))

	return service, nil
}

func (s *stateService) GetBalance() decimal.Decimal {
	s.stateM.RLock()
	defer s.stateM.RUnlock()
	return s.state.BudgetBalance
}

func (s *stateService) DecreaseBudgetBalance(v decimal.Decimal) (decimal.Decimal, error) {
	s.stateM.Lock()
	defer s.stateM.Unlock()
	s.state.BudgetBalance = s.state.BudgetBalance.Sub(v)
	if err := s.stateStorage.UpdateBalance(s.state.BudgetBalance); err != nil {
		return decimal.Decimal{}, err
	}
	return s.state.BudgetBalance, nil
}

func (s *stateService) DecreaseBudgetBalanceTx(tx *sqlx.Tx, v decimal.Decimal) (decimal.Decimal, error) {
	s.stateM.Lock()
	defer s.stateM.Unlock()
	s.state.BudgetBalance = s.state.BudgetBalance.Sub(v)
	log.Println("start UpdateBalanceTx")
	if err := s.stateStorage.UpdateBalanceTx(tx, s.state.BudgetBalance); err != nil {
		return decimal.Decimal{}, err
	}
	log.Println("end UpdateBalanceTx")
	return s.state.BudgetBalance, nil
}

func (s *stateService) runJob(ctx context.Context, nextTriggerTime time.Duration) {
	timer := time.NewTimer(nextTriggerTime)

	for {
		select {
		case <-timer.C:
			nextMonth := time.Now().AddDate(0, 1, 0).Truncate(24 * time.Hour)
			timer = time.NewTimer(nextMonth.Sub(time.Now()))

			if err := s.updateBalance(nextMonth); err != nil {
				log.Printf("error on update state, %v\n", err)
			}
		case <-ctx.Done():
			log.Printf("cancel update state job")
			return
		}
	}
}

func (s *stateService) updateBalance(nextTime time.Time) error {
	s.stateM.RLock()
	defer s.stateM.RUnlock()
	if err := s.stateStorage.UpdateBalanceAndExpiresIn(s.state.BudgetValue, nextTime); err != nil {
		return err
	}
	s.state.BudgetBalance = s.state.BudgetValue
	s.state.BudgetExpiresIn = nextTime

	return nil
}
