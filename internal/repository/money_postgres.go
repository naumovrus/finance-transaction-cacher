package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
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
	moneyTable              = "money_users"
	transactionSendTable    = "transactionssend"
	transactionBalanceTable = "transactionsbalance"
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
func (r *MoneyPostgres) TopUp(userId int, amount float64) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	query := fmt.Sprintf("UPDATE %s mt SET amount=amount+$1 WHERE mt.user_id = $2", moneyTable)
	_, err = r.db.Exec(query, amount, userId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	var id int
	queryCreateTr := fmt.Sprintf(`INSERT INTO %s (user_id, date_time) VALUES ($1, $2) RETURNING id`, transactionBalanceTable)
	row := r.db.QueryRow(queryCreateTr, userId, time.Now())
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}
	return id, tx.Commit()
}

// update logic
func (r *MoneyPostgres) TakeOut(userId int, amount float64) (int, error) {
	// add error when amount take out < money.amount

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	var amCheck float64
	queryCheck := fmt.Sprintf("SELECT amount FROM %s WHERE user_id = $1", moneyTable)

	err = r.db.Get(&amCheck, queryCheck, userId)
	if err != nil {
		tx.Rollback()
		return 0, nil
	}
	if amCheck < amount {
		tx.Rollback()
		return 0, errors.New("can't afford")
	}

	query := fmt.Sprintf("UPDATE %s mt SET amount=amount-$1 WHERE mt.user_id = $2", moneyTable)

	_, err = r.db.Exec(query, amount, userId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	var id int
	queryCreateTr := fmt.Sprintf(`INSERT INTO %s (user_id, date_time) VALUES ($1, $2) RETURNING id`, transactionBalanceTable)
	row := r.db.QueryRow(queryCreateTr, userId, time.Now())
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *MoneyPostgres) Send(uuid string, userIdFrom, userIdTo int, amount float64, time time.Time) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// var id int
	// add error when amount take out < money.amount
	var amCheck float64
	queryCheck := fmt.Sprintf("SELECT amount FROM %s WHERE user_id = $1", moneyTable)

	err = r.db.Get(&amCheck, queryCheck, userIdFrom)
	if err != nil {
		tx.Rollback()
		return err
	}
	if amCheck < amount {
		tx.Rollback()
		return errors.New("can't afford")
	}

	queryUserFrom := fmt.Sprintf(`UPDATE %s mt SET amount=amount-$1 WHERE mt.user_id = $2`, moneyTable)

	_, err = tx.Exec(queryUserFrom, amount, userIdFrom)

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
	// var data ent.TransactionSend

	queryCreateTr := fmt.Sprintf(`INSERT INTO %s (uuid, user_id_from, user_id_to, amount, date_time) VALUES ($1, $2, $3, $4, $5)`, transactionSendTable)
	r.db.QueryRow(queryCreateTr, uuid, userIdFrom, userIdTo, amount, time)
	tx.Commit()
	return nil

}

func (r *MoneyPostgres) GetLastTransactionSend() (int, error) {
	query := fmt.Sprintf("SELECT id FROM %s ORDER BY id DESC LIMIT 1", transactionSendTable)
	var id int
	err := r.db.Get(&id, query)
	switch err {
	case sql.ErrNoRows:
		return 1, nil
	case nil:
		return id, nil
	default:
		return 0, err
	}

}

func (r *MoneyPostgres) GetLastTransactionTUTO() (int, error) {
	query := fmt.Sprintf("SELECT id FROM %s ORDER BY desc LIMIT 1", transactionBalanceTable)
	var id int
	err := r.db.Get(&id, query)
	if err != nil {
		return 0, nil
	}
	if id == 0 {
		return 1, nil
	}
	return id, nil
}

func (r *MoneyPostgres) SetCachedDataSendPostgres(userIdFrom, userIdTo int, time time.Time) (int, error) {
	var id int
	// queryCheck := fmt.Sprintf("SELECT id FROM %s WHERE date_time=$1", transactionSendTable)
	// row := r.db.QueryRow(queryCheck, time)
	// if err := row.Scan(&id); err == nil {
	// 	return 0, err
	// }
	query := fmt.Sprintf("INSERT INTO %s (user_id_from, user_id_to, date_time) VALUES ($1, $2, $3) RETURNING id", transactionSendTable)
	row := r.db.QueryRow(query, userIdFrom, userIdTo, time)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil

}
