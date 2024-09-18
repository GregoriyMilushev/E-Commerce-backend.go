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
    orderController := controllers.NewOrderController(db)
    productController := controllers.NewProductController(db)

    api := router.Group("/api")
    {
        api.POST("/orders", orderController.CreateOrder)
        api.GET("/products", productController.GetProducts)
        api.GET("/products/:id", productController.ShowProductr)

        logged := api.Group("/",)
        logged.Use(middleware.RequireAuth)
        {
            logged.POST("/me", authController.Me)
            // Orders
            logged.GET("/orders/:id", orderController.ShowOrder)
            logged.GET("/orders", orderController.GetOrders)
            
            admin := logged.Group("/")
            admin.Use(middleware.RequireAdminRole)  
            {
                // Users
                admin.GET("/users", userController.GetUsers)    
                admin.POST("/users", userController.CreateUser) 
                // Orders
                admin.GET("/all-orders", orderController.GetAllOrders)
                admin.DELETE("/orders/:id", orderController.DeleteOrder)
                // Products
                admin.POST("/products", productController.CreateProduct)
                admin.PATCH("/products/:id", productController.UpdateProduct)
                admin.DELETE("/products/:id", productController.Delete)
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
