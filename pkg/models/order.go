package models

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	Price     float64 `gorm:"column:price; not null"`
	ProductID uint    `gorm:"column:product_id; not null"`
	UserID    uint    `gorm:"column:user_id; not null"`
	Quantity  uint    `gorm:"column:quantity; not null"`
}
