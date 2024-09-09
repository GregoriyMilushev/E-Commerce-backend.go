package controllers

import (
	"net/http"
	"pharmacy-backend/models"
	"pharmacy-backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
    userService *services.UserService
}

func NewUserController(db *gorm.DB) *UserController {
    return &UserController{
        userService: services.NewUserService(db),
    }
}

func (ctrl *UserController) GetUsers(c *gin.Context) {
    users, err := ctrl.userService.GetAllUsers()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, users)
}

func (ctrl *UserController) CreateUser(c *gin.Context) {
    
    var user models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if err := ctrl.userService.CreateUser(&user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, user)
}
