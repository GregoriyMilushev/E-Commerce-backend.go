package controllers

import (
	"net/http"
	"pharmacy-backend/models"
	"pharmacy-backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductController struct {
    productService *services.ProductService
}

func NewProductController(db *gorm.DB) *ProductController {
    return &ProductController{
        productService: services.NewProductService(db),
    }
}

func (pc *ProductController) GetProducts(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

    name := c.Query("name")
    minPrice, _ := strconv.ParseFloat(c.DefaultQuery("minPrice", "0"), 64)
    maxPrice, _ := strconv.ParseFloat(c.Query("maxPrice"), 64)

    products, totalItems, err := pc.productService.GetPaginatedProducts(page, pageSize, name, minPrice, maxPrice)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products", "details": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "products":     products,
        "totalItems":   totalItems,
        "page":         page,
        "pageSize":     pageSize,
        "totalPages":   (totalItems + int64(pageSize) - 1) / int64(pageSize), // Calculating total pages
    })
}


func (pc *ProductController) ShowProductr(c *gin.Context) {
    orderId, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
        return
    }

    order, err  := pc.productService.GetProductByID(uint(orderId))
    if  err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "message": "Order not found",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "order": order,
    })
}

func (pc *ProductController) CreateProduct(c *gin.Context) {
    var product models.Product
    if err := c.ShouldBindJSON(&product); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := pc.productService.CreateProduct(&product); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
        return
    }

    c.JSON(http.StatusCreated, product)
}

func (pc *ProductController) UpdateProduct(c *gin.Context) {
    productId, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
        return
    }

    var updatedProduct models.Product
    if err := c.ShouldBindJSON(&updatedProduct); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product data", "details": err.Error()})
        return
    }

    updatedProductResponse, err := pc.productService.UpdateProduct(uint(productId), &updatedProduct)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product", "details": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "product": updatedProductResponse,
    })
}

func (pc *ProductController) Delete(c *gin.Context) {
    id := c.Param("id")
    productID, err := strconv.ParseUint(id, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
        return
    }

    err = pc.productService.DeletePrduct(uint(productID))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
