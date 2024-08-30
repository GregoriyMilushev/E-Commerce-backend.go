package routes

import (
	"pharmacy-backend/controllers"
	"pharmacy-backend/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
    // err := router.SetTrustedProxies([]string{"192.168.1.0/24", "10.0.0.0/8"})
    // if err != nil {
    //     panic(err)
    // }
    userController := controllers.NewUserController(db)
    authController := controllers.NewAuthController(db)

    api := router.Group("/api")
    {
        logged := api.Group("/",)
        logged.Use(middleware.RequireAuth)
        {
            logged.GET("/users", userController.GetUsers)
            logged.POST("/users", userController.CreateUser)
            logged.POST("/me", authController.Me)
        }

        auth := router.Group("/auth")
        {
            auth.POST("/login", authController.Login)
            auth.POST("/register", authController.Register)
            // auth.GET("/google/login", authController.GoogleLogin)
            // auth.GET("/google/callback", authController.GoogleCallback)
        }
    }
}
