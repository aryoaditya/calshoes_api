package models

import "time"

type Payment struct {
	Id        uint      `json:"id" gorm:"primaryKey"`
	OrderId   uint      `json:"order_id"`
	Method    string    `json:"method" gorm:"size:50"`
	Status    string    `json:"status" gorm:"size:50"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Order     Order     `gorm:"foreignKey:OrderId"`
}
