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
    // orderController := controllers.NewOrderController(db)

    api := router.Group("/api")
    {
        // api.GET("/orders/:id", orderController.GetOrderById)

        logged := api.Group("/",)
        logged.Use(middleware.RequireAuth)
        {
            logged.POST("/me", authController.Me)

            admin := logged.Group("/")
            admin.Use(middleware.RequireAdminRole)  
            {
                admin.GET("/users", userController.GetUsers)    
                admin.POST("/users", userController.CreateUser) 
            }
        }
    }

    auth := router.Group("/auth")
    {
        auth.POST("/login", authController.Login)
        auth.POST("/register", authController.Register)
        // auth.GET("/google/login", authController.GoogleLogin)
        // auth.GET("/google/callback", authController.GoogleCallback)
    }
}
