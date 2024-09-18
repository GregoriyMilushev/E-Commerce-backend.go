package services

import (
	"fmt"
	"pharmacy-backend/models"
	"time"

	"gorm.io/gorm"
)

type ProductResponse struct {
    ID            uint      `json:"id"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
    Name          string    `json:"name"`
    Description   string   `json:"description"`
    Price         float64 `json:"price"`
    Stock         uint `json:"stock"`
    Brand         models.Brand `json:"brand"`
}

type ProductService struct {
    db *gorm.DB
}

func NewProductService(db *gorm.DB) *ProductService {
    return &ProductService{db: db}
}

func (ps *ProductService) GetPaginatedProducts(page int, pageSize int, name string, minPrice float64, maxPrice float64) ([]models.Product, int64, error) {
    var products []models.Product
    var totalItems int64

    query := ps.db.Model(&models.Product{})

    if name != "" {
        query = query.Where("name LIKE ?", "%"+name+"%")
    }
    if minPrice > 0 {
        query = query.Where("price >= ?", minPrice)
    }
    if maxPrice > 0 {
        query = query.Where("price <= ?", maxPrice)
    }

    if err := query.Count(&totalItems).Error; err != nil {
        return nil, 0, err
    }

    offset := (page - 1) * pageSize
    if err := query.Limit(pageSize).Offset(offset).Preload("Brand").Find(&products).Error; err != nil {
        return nil, 0, err
    }

    return products, totalItems, nil
}


func (s *ProductService) GetProductByID(id uint) (*ProductResponse, error) {
    var product models.Product
    if err := s.db.Preload("Brand").First(&product, id).Error; err != nil {
        return nil, err
    }

    productResponse := &ProductResponse{
        ID:        product.ID,
        CreatedAt: product.CreatedAt,
        UpdatedAt: product.UpdatedAt,
        Name:    product.Name,
        Price:     product.Price,
        Stock:     product.Stock,
        Description: product.Description,
        Brand: product.Brand,
    }

   

    return productResponse, nil
}

func (ps *ProductService) CreateProduct(product *models.Product) error {
    if err := ps.db.Create(product).Error; err != nil {
        return err
    }

    return  nil
}

func (ps *ProductService) UpdateProduct(id uint, updatedProduct *models.Product) (*ProductResponse, error) {
    var product models.Product
    if err := ps.db.Model(&product).Where("id = ?", id).Updates(updatedProduct).
        Preload("Brand").First(&product).Error; err != nil {
        return nil, fmt.Errorf("failed to update product with id %d: %w", id, err)
    }

    productResponse := &ProductResponse{
        ID:          product.ID,
        CreatedAt:   product.CreatedAt,
        UpdatedAt:   product.UpdatedAt,
        Name:        product.Name,
        Price:       product.Price,
        Stock:       product.Stock,
        Description: product.Description,
        Brand:       product.Brand,
    }

    return productResponse, nil
}

func (ps *ProductService) DeletePrduct(id uint) error {
    result := ps.db.Delete(&models.Product{}, id)
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return fmt.Errorf("product with ID %d not found", id)
    }
    return nil
}


