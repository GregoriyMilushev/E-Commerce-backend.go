package middleware

import (
	"fmt"
	"net/http"
	"pharmacy-backend/config"
	"pharmacy-backend/database"
	"pharmacy-backend/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *gin.Context) {
	tokenString, err := c.Cookie("Authoriztion")
	if err != nil {
		fmt.Println("1111111")
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("222")
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
	
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(config.LoadConfig().JwtSectet), nil
	})
	
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			fmt.Println("333")
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		var user models.User
	 	if userError :=	database.GetDB().First(&user, claims["user_id"]).Error; userError != nil {
			fmt.Println("444")
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		c.Set("user", &user)
		c.Next()
	} else {
		fmt.Println("666")
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}