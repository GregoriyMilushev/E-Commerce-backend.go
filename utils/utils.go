package utils

import (
	"fmt"
	"pharmacy-backend/config"
	"pharmacy-backend/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)



func GenerateJWT(user models.User) (string, error) {
	jwtSecret := config.LoadConfig().JwtSectet
	if jwtSecret == "" {
        return "", fmt.Errorf("JWT_SECRET is not set in the environment")
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "email":   user.Email,
        "exp":     time.Now().Add(time.Hour * 24 * 30).Unix(), // Token expires in 24 hours
    })

    tokenString, err := token.SignedString([]byte(jwtSecret))
    if err != nil {
        return "", err
    }

    return tokenString, nil
}