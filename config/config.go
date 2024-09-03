package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
    TZ                   string
    GoEnv                string
    NotificationMode     string
    AppPort              string
    AppName              string
    ApiPrefix            string
    AppFallbackLanguage  string
    AppHeaderLanguage    string
    FrontendDomain       string
    BackendDomain        string
    DatabaseType         string
    DatabaseHost         string
    DatabasePort         string
    DatabaseUsername     string
    DatabasePassword     string
    DatabaseName         string
    JwtSectet            string
}

var (
	config   *Config
	once sync.Once
)

func LoadConfig() *Config {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }
    once.Do(func() {
        config = &Config{
            TZ:                   getEnv("TZ", "UTC"),
            GoEnv:                getEnv("GO_ENV", "development"),
            NotificationMode:     getEnv("NOTIFICATION_MODE", "development"),
            AppPort:              getEnv("APP_PORT", "8888"),
            AppName:              getEnv("APP_NAME", "pharmacy-backend"),
            ApiPrefix:            getEnv("API_PREFIX", "api/:client"),
            AppFallbackLanguage:  getEnv("APP_FALLBACK_LANGUAGE", "en"),
            AppHeaderLanguage:    getEnv("APP_HEADER_LANGUAGE", "x-custom-lang"),
            FrontendDomain:       getEnv("FRONTEND_DOMAIN", "http://localhost:3011"),
            BackendDomain:        getEnv("BACKEND_DOMAIN", "http://localhost:8888"),
            DatabaseType:         getEnv("DATABASE_TYPE", "mysql"),
            DatabaseHost:         getEnv("DATABASE_HOST", "localhost"),
            DatabasePort:         getEnv("DATABASE_PORT", "3306"),
            DatabaseUsername:     getEnv("DATABASE_USERNAME", "root"),
            DatabasePassword:     getEnv("DATABASE_PASSWORD", "root"),
            DatabaseName:         getEnv("DATABASE_NAME", "pharmacy-test"),
            JwtSectet:            getEnv("JWT_SECRET", ""),
        }
	})


    return config
}

func getEnv(key, fallback string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }

    return fallback
}

func ConnectDatabase(databaseURL string) (*gorm.DB, error) {
    db, err := gorm.Open(mysql.Open(databaseURL), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    return db, nil
}

func GetDatabaseURL(config *Config) string {
    return config.DatabaseUsername + ":" + config.DatabasePassword + "@tcp(" + config.DatabaseHost + ":" + config.DatabasePort + ")/" + config.DatabaseName + "?parseTime=True"
}
