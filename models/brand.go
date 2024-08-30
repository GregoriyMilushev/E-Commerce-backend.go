package models

import (
	"gorm.io/gorm"
)


type Brand struct {
    gorm.Model
    Name        string  `gorm:"size:100;not null" json:"name"`
    Description string  `gorm:"type:text" json:"description"`
    Discount    float64 `gorm:"type:decimal(10,2)" json:"discount"`
    Products    []Product `gorm:"foreignKey:BrandID"`
}
