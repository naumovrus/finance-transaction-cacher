package service

import (
	"time"

	"github.com/naumovrus/finance-transaction-api/internal/repository"
)

type MoneySerice struct {
	repo repository.Money
}

func NewMoneyService(repo repository.Money) *MoneySerice {
	return &MoneySerice{repo: repo}
}

func (s *MoneySerice) CreateWallet(userId int) (int, error) {
	return s.repo.CreateWallet(userId)
}

func (s *MoneySerice) TopUp(userId int, amount float64) (int, error) {
	return s.repo.TopUp(userId, amount)
}

func (s *MoneySerice) TakeOut(userId int, amount float64) (int, error) {
	return s.repo.TakeOut(userId, amount)
}

func (s *MoneySerice) Send(userIdFrom, userIdTo int, amount float64) (int, error) {
	return s.repo.Send(userIdFrom, userIdTo, amount)
}

func (s *MoneySerice) GetLastTransactionSend() (int, error) {
	return s.repo.GetLastTransactionSend()
}

func (s *MoneySerice) GetLastTransactionTUTO() (int, error) {
	return s.repo.GetLastTransactionTUTO()
}

func (s *MoneySerice) SetCachedDataSendPostgres(userIdFrom, userIdTo int, time time.Time) (int, error) {
	return s.repo.SetCachedDataSendPostgres(userIdFrom, userIdTo, time)
}
