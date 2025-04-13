package models

import (
	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model
	UserID   uint      `json:"user_id"`
	Products []Product `gorm:"many2many:cart_products;" json:"products"`
}
