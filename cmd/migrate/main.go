package main

import (
	"database/sql"
	"log"
	"pharmacy-backend/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func runMigrations(db *sql.DB, dbType string) error {
    // Initialize the migrate instance
    driver, err := mysql.WithInstance(db, &mysql.Config{})
    if err != nil {
        return err
    }

    m, err := migrate.NewWithDatabaseInstance(
        "file://migrations",
        dbType, 
        driver)
    if err != nil {
        return err
    }
    
    // Apply Up migrations
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return err
    }

    if version, _, verisonErr := m.Version(); verisonErr != nil{
        return verisonErr
    } else if version > 0  {
        log.Printf("Migrations version: %v applied successfully!", version)
    } 

    return nil
}

func connectDatabase(dbType string,dbUrl string) (*sql.DB, error) {
    db, err := sql.Open(dbType, dbUrl)
    if err != nil {
        return nil, err
    }
    return db, nil
}

func main() {
    conf := config.LoadConfig()

    db, err := connectDatabase(conf.DatabaseType, config.GetDatabaseURL(conf))
    if err != nil {
        log.Fatalf("Could not connect to the database: %v", err)
    }
    defer db.Close()

    if err := runMigrations(db, conf.DatabaseType); err != nil {
        log.Fatalf("Could not apply migrations: %v", err)
    }
}
