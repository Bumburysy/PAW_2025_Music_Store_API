package middleware

import (
	"music-store-api/config"
	"music-store-api/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Brak nagłówka Authorization"})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := config.ValidateJWT(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Nieprawidłowy token"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)
		c.Next()
	}
}

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Brak danych roli w tokenie"})
			c.Abort()
			return
		}

		userRole := roleVal.(string)
		for _, role := range allowedRoles {
			if userRole == role {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Brak uprawnień do wykonania tej operacji"})
		c.Abort()
	}
}
