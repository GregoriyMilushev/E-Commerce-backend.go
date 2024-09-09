package services

import (
	"pharmacy-backend/models"

	"gorm.io/gorm"
)

type OrderService struct {
    db *gorm.DB
}

func NewOrderService(db *gorm.DB) *OrderService {
    return &OrderService{db: db}
}

func (s *OrderService) GetAllOrders() ([]models.Order, error) {
    var orders []models.Order
    if err := s.db.Find(&orders).Error; err != nil {
        return nil, err
    }
    return orders, nil
}


