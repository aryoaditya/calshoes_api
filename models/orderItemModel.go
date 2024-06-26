package models

import "time"

type OrderItem struct {
	Id        uint      `gorm:"primaryKey"`
	OrderId   uint      `json:"order_id"`
	ProductId uint      `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Product   Product   `gorm:"foreignKey:ProductId"`
	Order     Order     `gorm:"foreignKey:OrderId"`
}
