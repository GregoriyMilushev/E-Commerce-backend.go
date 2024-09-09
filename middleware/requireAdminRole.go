package middleware

import (
	"net/http"
	"pharmacy-backend/models"

	"github.com/gin-gonic/gin"
)

func RequireAdminRole(c *gin.Context) {
    user, exists := c.Get("user")
    if !exists || user.(*models.User).Role != models.RoleAdmin {
        c.AbortWithStatus(http.StatusForbidden)
        return
    }

    c.Next()
}
