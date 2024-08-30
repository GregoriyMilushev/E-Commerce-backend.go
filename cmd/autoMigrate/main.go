package main

import (
	"log"
	"pharmacy-backend/config"
	"pharmacy-backend/models"
)

func main() {
    log.Println("Starting the application...")

    conf := config.LoadConfig()

    db, err := config.ConnectDatabase(config.GetDatabaseURL(conf))
    if err != nil {
        log.Fatalf("Could not connect to the database: %v", err)
    }

    // Auto migrate models you want by adding them here
    if err := db.AutoMigrate(&models.User{}, &models.Product{}); err != nil {
        log.Fatalf("Could not migrate database: %v", err)
    }

    log.Println("Database connected and models migrated successfully!")
}
