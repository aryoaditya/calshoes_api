package models

import "time"

type Order struct {
	Id         uint      `json:"id" gorm:"primaryKey"`
	CustomerId uint      `json:"customer_id"`
	TotalPrice float64   `json:"total_price"`
	Status     string    `json:"status" gorm:"size:50"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Customer   Customer  `gorm:"foreignKey:CustomerId"`
	OrderItems []OrderItem
}
