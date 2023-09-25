package service

import (
	ent "github.com/naumovrus/finance-transaction-api/internal/entity"
	"github.com/naumovrus/finance-transaction-api/internal/repository"
)

type User interface {
	CreateUser(user ent.User) (int, error)
	GetDataUser(userId int) (ent.User, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

// update
type Money interface {
	CreateWallet(userId int) (int, error)
	TopUp(userId int, amount float64) error
	TakeOut(userId int, amount float64) error
	Send(userIdFrom, userIdTo int, amount float64) error
}

type Service struct {
	User
	Money
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		User:  NewUserService(repo.User),
		Money: NewMoneyService(repo.Money),
	}

}
