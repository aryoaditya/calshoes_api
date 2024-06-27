package models

import "time"

type Customer struct {
	Id        uint      `json:"id" gorm:"primaryKey"`
	FirstName string    `json:"first_name" gorm:"size:50"`
	LastName  *string   `json:"last_name" gorm:"size:100"`
	Email     string    `gorm:"unique; size:100"`
	Password  string    `json:"password" gorm:"size:100"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Carts     []Cart
	Orders    []Order
}
