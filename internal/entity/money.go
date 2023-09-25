package entity

type Money struct {
	Id     int     `json:"id" db:"id"`
	UserId int     `json:"user_id" binging:"required"`
	Amount float64 `json:"amount" binding:"required"`
}
