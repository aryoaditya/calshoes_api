package models

import "time"

type Product struct {
	Id          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:100"`
	Description *string   `json:"description" gorm:"type:text"`
	Price       float64   `json:"price"`
	ImargeUrl   string    `json:"image_url" gorm:"size:100"`
	CategoryId  uint      `json:"category_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at" gorm:"index"`
	Category    Category  `gorm:"foreignKey:CategoryId"`
}
