package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type MoneyPostgres struct {
	db *sqlx.DB
}

func NewMoneyPostgres(db *sqlx.DB) *MoneyPostgres {
	return &MoneyPostgres{
		db: db,
	}
}

const (
	moneyTable       = "money_users"
	transactionTable = "transactions"
)

func (r *MoneyPostgres) CreateWallet(userId int) (int, error) {
	var id int

	createListQuery := fmt.Sprintf("INSERT INTO %s (user_id) VALUES ($1) RETURNING id", moneyTable)
	row := r.db.QueryRow(createListQuery, userId)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil

}

// update logic
func (r *MoneyPostgres) TopUp(userId int, amount float64) error {
	query := fmt.Sprintf("UPDATE %s mt SET amount=amount+$1 WHERE mt.user_id = $2", moneyTable)
	_, err := r.db.Exec(query, amount, userId)
	if err != nil {
		return err
	}
	return nil
}

// update logic
func (r *MoneyPostgres) TakeOut(userId int, amount float64) error {
	// add error when amount take out < money.amount
	var amCheck float64
	queryCheck := fmt.Sprintf("SELECT amount FROM %s WHERE user_id = $1", moneyTable)

	err := r.db.Get(&amCheck, queryCheck, userId)

	// ok = amCheck < amount

	query := fmt.Sprintf("UPDATE %s mt SET amount=amount-$1 WHERE mt.user_id = $2", moneyTable)

	_, err = r.db.Exec(query, amount, userId)
	if err != nil {
		return err
	}
	return nil
}

func (r *MoneyPostgres) Send(userIdFrom, userIdTo int, amount float64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// var id int
	// add error when amount take out < money.amount
	queryUserFrom := fmt.Sprintf(`UPDATE %s mt SET amount=amount-$1 WHERE mt.user_id = $2`, moneyTable)

	_, err = tx.Exec(queryUserFrom, amount, userIdFrom)
	logrus.Printf("transaction 1 succsess")
	if err != nil {
		tx.Rollback()
		return err
	}

	queryUserTo := fmt.Sprintf(`UPDATE %s mt SET amount=amount+$1 WHERE mt.user_id=$2`, moneyTable)

	_, err = tx.Exec(queryUserTo, amount, userIdTo)
	if err != nil {
		tx.Rollback()
		return err
	}
	logrus.Printf("")
	return tx.Commit()

	// queryCreateTr := fmt.Sprintf(`INSERT INTO %s (user_id_from, user_id_to, time) VALUES ($1, $2, $3) RETURINING id`, transactionTable)
	// row := r.db.QueryRow(queryCreateTr, userIdFrom, userIdTo, time.Now())
	// if err := row.Scan(&id); err != nil {
	// 	return 0, err
	// }
	// return id, nil
}
