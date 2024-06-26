package models

import "time"

type CartItem struct {
	Id        uint      `json:"id" gorm:"primaryKey"`
	CartId    uint      `json:"cart_id"`
	ProductId uint      `json:"product_id"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Cart      Cart      `gorm:"foreignKey:CartId"`
	Product   Product   `gorm:"foreignKey:ProductId"`
}
