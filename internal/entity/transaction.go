package entity

import "time"

type TransactionSend struct {
	Id         int       `json:"id" db:"id"`
	UserIdFrom int       `json:"user_id_from" binding:"required"`
	UserIdTo   int       `json:"user_id_to" binding:"required"`
	Time       time.Time `json:"time" binding:"required"`
}

type TransactionTUTO struct {
	Id     int       `json:"-" db:"id"`
	UserId int       `json:"user_id" bindning:"required"`
	Time   time.Time `json:"time" binding:"required"`
}
