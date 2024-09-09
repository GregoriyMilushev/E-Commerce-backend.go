package controllers

import (
	"net/http"
	"pharmacy-backend/services"

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


func (oc *OrderController) GetOrders(c *gin.Context) {
    orders, err := oc.orderService.GetAllOrders()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, orders)
}

// func (oc *OrderController) ShowOrder(c *gin.Context) {
//     orderId := c.Param("id")

//     var order models.Order
//     if err := oc.db.Where("id = ?", orderId).Preload("Products").First(&order).Error; err != nil {
//         // If the order is not found, return a 404 error
//         c.JSON(http.StatusNotFound, gin.H{
//             "message": "Order not found",
//         })
//         return
//     }

//     // If the order is found, return it
//     c.JSON(http.StatusOK, gin.H{
//         "order": order,
//     })
// }