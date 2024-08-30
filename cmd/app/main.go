package main

import (
	"log"
	"pharmacy-backend/config"
	"pharmacy-backend/database"
	"pharmacy-backend/routes"

	"github.com/gin-gonic/gin"
)

// func init() {

// }

func main() {
    conf := config.LoadConfig()

    router := gin.Default()

    // Initialize the Database Once
    db := database.GetDB()

    // Setup routes
    routes.SetupRoutes(router, db)

    // Start the server
    if err := router.Run(":" + conf.AppPort); err != nil {
        log.Fatalf("Could not start server: %v", err)
    }
}
