package entity

import "time"

type TransactionSend struct {
	Uuid       string    `json:"uuid" db:"uuid"`
	UserIdFrom int       `json:"user_id_from" binding:"required"`
	UserIdTo   int       `json:"user_id_to" binding:"required"`
	Amount     float64   `json:"amount" binding:"required"`
	Time       time.Time `json:"time" binding:"required"`
}

type TransactionTUTO struct {
	Uuid   string    `json:"uuid" db:"uuid"`
	UserId int       `json:"user_id" bindning:"required"`
	Time   time.Time `json:"time" binding:"required"`
}
