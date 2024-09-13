package services

import (
	"fmt"
	"log"
	"pharmacy-backend/models"
	"time"

	"gorm.io/gorm"
)

type OrderResponse struct {
    ID            uint      `json:"id"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
    UserID        uint      `json:"user_id"`
    Total         float64   `json:"total"`
    OrderProducts []OrderProductResponse `json:"order_products"`
}

type OrderProductResponse struct {
    ID        uint    `json:"id"`
    ProductID uint    `json:"product_id"`
    Quantity  uint     `json:"quantity"`
    Price     float64 `json:"price"`
}

type OrderService struct {
    db *gorm.DB
}

func NewOrderService(db *gorm.DB) *OrderService {
    return &OrderService{db: db}
}

func (s *OrderService) GetAllOrders() ([]models.Order, error) {
    var orders []models.Order
    if err := s.db.Preload("OrderProducts.Product").Find(&orders).Error; err != nil {
        return nil, err
    }
    return orders, nil
}

func (s *OrderService) GetOrders(userID uint) ([]models.Order, error) {
    var orders []models.Order
    if err := s.db.Preload("OrderProducts.Product").Where("user_id = ?",userID).Find(&orders).Error; err != nil {
        return nil, err
    }
    return orders, nil
}

func (s *OrderService) GetOrderByID(id uint) (*OrderResponse, error) {
    var order models.Order
    if err := s.db.Preload("OrderProducts.Product").First(&order, id).Error; err != nil {
        return nil, err
    }

    orderResponse := &OrderResponse{
        ID:        order.ID,
        CreatedAt: order.CreatedAt,
        UpdatedAt: order.UpdatedAt,
        UserID:    order.UserID,
        Total:     order.Total,
    }

    for _, op := range order.OrderProducts {
        orderResponse.OrderProducts = append(orderResponse.OrderProducts, OrderProductResponse{
            ID:        op.ID,
            ProductID: op.ProductID,
            Quantity:  op.Quantity,
            Price:     op.Price,
        })
    }

    return orderResponse, nil
}

func (os *OrderService) CreateOrder(userID uint, orderProducts []models.OrderProduct) (*OrderResponse, error) {
    order := &models.Order{
        UserID:        userID,
        Total: 0,
    }

    err := os.db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Create(order).Error; err != nil {
            return err
        }
        
        for i := range orderProducts {
            var product models.Product
            if err := tx.First(&product, orderProducts[i].ProductID).Error; err != nil {
                return err

            }

            if product.Stock < uint(orderProducts[i].Quantity) {
                return fmt.Errorf("insufficient stock for product %d", orderProducts[i].ProductID)

            }

            product.Stock -= uint(orderProducts[i].Quantity)
            if err := tx.Save(&product).Error; err != nil {
                return err
            }

            log.Printf("Product Price - %v", product.Price)
            orderProducts[i].Price = product.Price * float64(orderProducts[i].Quantity)
            orderProducts[i].OrderID = order.ID
            order.Total += orderProducts[i].Price
        }
        order.OrderProducts = orderProducts
        if err := tx.Save(&order).Error; err != nil {
            return err
        }

        return nil
    })

    if err != nil {
        return nil, err
    }

    orderResponse := &OrderResponse{
        ID:        order.ID,
        CreatedAt: order.CreatedAt,
        UpdatedAt: order.UpdatedAt,
        UserID:    order.UserID,
        Total:     order.Total,
    }

    for _, op := range order.OrderProducts {
        orderResponse.OrderProducts = append(orderResponse.OrderProducts, OrderProductResponse{
            ID:        op.ID,
            ProductID: op.ProductID,
            Quantity:  op.Quantity,
            Price:     op.Price,
        })
    }

    return orderResponse, nil
}

func (s *OrderService) DeleteOrder(id uint) error {
    result := s.db.Delete(&models.Order{}, id)
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return fmt.Errorf("order with ID %d not found", id)
    }
    return nil
}
