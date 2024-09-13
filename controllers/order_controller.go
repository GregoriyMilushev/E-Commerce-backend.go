package controllers

import (
	"net/http"
	"pharmacy-backend/models"
	"pharmacy-backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


type OrderController struct {
    orderService *services.OrderService
}

func NewOrderController(db *gorm.DB) *OrderController {
    return &OrderController{
        orderService: services.NewOrderService(db),
    }
}


func (oc *OrderController) GetAllOrders(c *gin.Context) {
    orders, err := oc.orderService.GetAllOrders()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, orders)
}

func (oc *OrderController) GetOrders(c *gin.Context) {
    user, _ := c.Get("user")
    orders, err := oc.orderService.GetOrders(user.(*models.User).ID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, orders)
}

func (oc *OrderController) ShowOrder(c *gin.Context) {
    orderId, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
        return
    }

    order, err  := oc.orderService.GetOrderByID(uint(orderId))
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

func (oc *OrderController) CreateOrder(c *gin.Context) {
    var orderRequest struct {
        UserID    uint `json:"userId" binding:"required"`
        Products  []struct {
            ProductID uint `json:"productId" binding:"required"`
            Quantity  uint  `json:"quantity" binding:"required"`
        } `json:"products"`
    }
    
    if err := c.ShouldBindJSON(&orderRequest); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    orderProducts := make([]models.OrderProduct, len(orderRequest.Products))
    for i, p := range orderRequest.Products {
        orderProducts[i] = models.OrderProduct{
            ProductID: p.ProductID,
            Quantity:  p.Quantity,
        }
    }
    order, err := oc.orderService.CreateOrder(orderRequest.UserID, orderProducts)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"order": order})
}

func (oc *OrderController) DeleteOrder(c *gin.Context) {
    id := c.Param("id")
    orderID, err := strconv.ParseUint(id, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
        return
    }

    err = oc.orderService.DeleteOrder(uint(orderID))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Order soft deleted successfully"})
}
