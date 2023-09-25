package repository

import (
	"github.com/jmoiron/sqlx"
	ent "github.com/naumovrus/finance-transaction-api/internal/entity"
)

type User interface {
	CreateUser(user ent.User) (int, error)
	GetUser(username, password string) (ent.User, error)
	GetDataUser(userId int) (ent.User, error)
}

// update
type Money interface {
	CreateWallet(userId int) (int, error)
	TopUp(userId int, amount float64) error
	TakeOut(userId int, amount float64) error
	Send(userIdFrom, userIdTo int, amount float64) error
}

type Repository struct {
	User
	Money
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User:  NewUserPostgres(db),
		Money: NewMoneyPostgres(db),
	}

}
