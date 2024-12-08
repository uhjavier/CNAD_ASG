package middleware

import (
	"car-sharing-system/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		// For demo purposes, create a dummy user
		// In production, you would validate the token and get real user data
		dummyUser := &models.User{
			Model: gorm.Model{ID: 1},
			Email: "user1@example.com",
			MembershipType: "BASIC",
		}
		
		c.Set("user", dummyUser)
		c.Next()
	}
} 