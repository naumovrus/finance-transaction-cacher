package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	ent "github.com/naumovrus/finance-transaction-api/internal/entity"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{
		db: db,
	}
}

const (
	usersTable = "users"
)

func (r *UserPostgres) CreateUser(user ent.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (first_name, last_name, username, password_hash) values ($1, $2, $3, $4) RETURNING id", usersTable)
	row := r.db.QueryRow(query, user.FirstName, user.LastName, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *UserPostgres) GetUser(username, password string) (ent.User, error) {
	var user ent.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password_hash=$2", usersTable)
	err := r.db.Get(&user, query, username, password)
	return user, err
}

func (r *UserPostgres) GetDataUser(userId int) (ent.User, error) {
	var user ent.User
	query := fmt.Sprintf("SELECT first_name, last_name, username, amount FROM %s WHERE id = $1", usersTable)
	err := r.db.Get(&user, query, userId)
	return user, err
}
