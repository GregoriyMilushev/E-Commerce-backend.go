package database

import (
	"log"
	"pharmacy-backend/config"
	"sync"

	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

func GetDB() *gorm.DB {
	once.Do(func() {
		conf := config.LoadConfig()
		var err error
		db, err = config.ConnectDatabase(config.GetDatabaseURL(conf))
		if err != nil {
			log.Fatalf("Could not connect to the database: %v", err)
		}
	})
	
	return db
}
