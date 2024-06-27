package models

import "time"

type Cart struct {
	Id         uint      `json:"id" gorm:"primaryKey"`
	CustomerId uint      `json:"customer_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Customer   Customer  `gorm:"foreignKey:CustomerId"`
}
