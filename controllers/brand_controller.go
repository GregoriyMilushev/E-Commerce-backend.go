package controllers

import (
	"net/http"
	"pharmacy-backend/models"
	"pharmacy-backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BrandController struct {
    brandService *services.BrandService
}

func NewBrandController(db *gorm.DB) *BrandController {
    return &BrandController{
        brandService: services.NewBrandService(db),
    }
}

func (pc *BrandController) GetBrands(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
    name := c.Query("name")

    brands, totalItems, err := pc.brandService.GetPaginatedBrands(page, pageSize, name)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve brands", "details": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "items":     brands,
        "totalItems":   totalItems,
        "page":         page,
        "pageSize":     pageSize,
        "totalPages":   (totalItems + int64(pageSize) - 1) / int64(pageSize), // Calculating total pages
    })
}


func (pc *BrandController) ShowBrand(c *gin.Context) {
    brandId, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid brand ID"})
        return
    }

    brand, err  := pc.brandService.GetBrandByID(uint(brandId))
    if  err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "message": "Brand not found",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "brand": brand,
    })
}

func (pc *BrandController) CreateBrand(c *gin.Context) {
    var brand models.Brand
    if err := c.ShouldBindJSON(&brand); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := pc.brandService.CreateBrand(&brand); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create brand"})
        return
    }

    c.JSON(http.StatusCreated, brand)
}

func (pc *BrandController) UpdateBrand(c *gin.Context) {
    brandId, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid brand ID"})
        return
    }

    var updatedBrand models.Brand
    if err := c.ShouldBindJSON(&updatedBrand); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid brand data", "details": err.Error()})
        return
    }

    updatedProductResponse, err := pc.brandService.UpdateBrand(uint(brandId), &updatedBrand)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update brand", "details": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "brand": updatedProductResponse,
    })
}

func (pc *BrandController) Delete(c *gin.Context) {
    id := c.Param("id")
    brandID, err := strconv.ParseUint(id, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
        return
    }

    err = pc.brandService.DeleteBrand(uint(brandID))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Brand deleted successfully"})
}
