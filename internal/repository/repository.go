package repository

import (
	"time"

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
	TopUp(userId int, amount float64) (int, error)
	TakeOut(userId int, amount float64) (int, error)
	Send(userIdFrom, userIdTo int, amount float64) (int, error)
	GetLastTransactionSend() (int, error)
	GetLastTransactionTUTO() (int, error)
	SetCachedDataSendPostgres(userIdFrom, userIdTo int, time time.Time) (int, error)
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
