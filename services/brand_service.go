package services

import (
	"fmt"
	"pharmacy-backend/models"
	"time"

	"gorm.io/gorm"
)

type BrandResponse struct {
    ID            uint      `json:"id"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
    Name          string    `json:"name"`
    Description   string   `json:"description"`
    Discount      float64 `json:"discount"`
    Products      []models.Product `json:"products"`
}

type BrandService struct {
    db *gorm.DB
}

func NewBrandService(db *gorm.DB) *BrandService {
    return &BrandService{db: db}
}

func (ps *BrandService) GetPaginatedBrands(page int, pageSize int, name string) ([]models.Brand, int64, error) {
    var brands []models.Brand
    var totalItems int64

    query := ps.db.Model(&models.Brand{})

    if name != "" {
        query = query.Where("name LIKE ?", "%"+name+"%")
    }

    if err := query.Count(&totalItems).Error; err != nil {
        return nil, 0, err
    }

    offset := (page - 1) * pageSize
    if err := query.Limit(pageSize).Offset(offset).Preload("Products").Find(&brands).Error; err != nil {
        return nil, 0, err
    }

    return brands, totalItems, nil
}


func (s *BrandService) GetBrandByID(id uint) (*BrandResponse, error) {
    var brand models.Brand
    if err := s.db.Preload("Products").First(&brand, id).Error; err != nil {
        return nil, err
    }

    brandResponse := &BrandResponse{
        ID:        brand.ID,
        CreatedAt: brand.CreatedAt,
        UpdatedAt: brand.UpdatedAt,
        Name:    brand.Name,
        Description: brand.Description,
        Discount: brand.Discount,
        Products: brand.Products,
    }

   

    return brandResponse, nil
}

func (ps *BrandService) CreateBrand(brand *models.Brand) error {
    if err := ps.db.Create(brand).Error; err != nil {
        return err
    }

    return  nil
}

func (ps *BrandService) UpdateBrand(id uint, updatedBrand *models.Brand) (*BrandResponse, error) {
    var brand models.Brand
    if err := ps.db.Model(&brand).Where("id = ?", id).Updates(updatedBrand).Preload("Products").First(&brand).Error; err != nil {
        return nil, fmt.Errorf("failed to update brand with id %d: %w", id, err)
    }

    brandResponse := &BrandResponse{
        ID:          brand.ID,
        CreatedAt:   brand.CreatedAt,
        UpdatedAt:   brand.UpdatedAt,
        Name:        brand.Name,
        Description: brand.Description,
        Discount:    brand.Discount,
        Products:    brand.Products,
    }

    return brandResponse, nil
}

func (ps *BrandService) DeleteBrand(id uint) error {
    result := ps.db.Delete(&models.Brand{}, id)
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return fmt.Errorf("brand with ID %d not found", id)
    }
    return nil
}


