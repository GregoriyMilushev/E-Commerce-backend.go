package models

import (
	"gorm.io/gorm"
)


type Product struct {
    gorm.Model
    Name        string  `gorm:"size:100;not null" json:"name"`
    Description string  `gorm:"type:text" json:"description"`
    Price       float64 `gorm:"type:decimal(10,2);not null" json:"price"`
    Discount    float64 `gorm:"type:decimal(10,2)" json:"discount"`
    Stock       uint    `gorm:"not null" json:"stock"`
    BrandID     uint    `json:"brand_id"`
    Brand       Brand   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"brand"`
    OrderProducts []OrderProduct  `gorm:"foreignKey:ProductID" json:"order_products"`
}
