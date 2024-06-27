package models

import "time"

type Category struct {
	Id        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:50"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
