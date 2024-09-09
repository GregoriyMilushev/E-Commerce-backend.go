package models

import (
	"gorm.io/gorm"
)

type Order struct {
    gorm.Model
    UserID    uint         `gorm:"not null" json:"user_id"`
    User      User         `json:"user"`
    OrderProducts []OrderProduct `gorm:"foreignKey:OrderID" json:"order_products"`
    Total     float64      `gorm:"type:decimal(10,2);not null" json:"total"`
}

type OrderProduct struct {
    ID        uint     `gorm:"primaryKey" json:"id"`
    OrderID   uint    `gorm:"not null" json:"order_id"`
    ProductID uint    `gorm:"not null" json:"product_id"`
    Quantity  uint     `gorm:"not null" json:"quantity"`
    Price     float64  `gorm:"type:decimal(10,2);not null" json:"price"`
 	Order     Order    `gorm:"foreignKey:OrderID" json:"order"`
    Product   Product  `gorm:"foreignKey:ProductID" json:"product"`
}

