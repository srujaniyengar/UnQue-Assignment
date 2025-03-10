package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"

	"UnQue/configs"
	"UnQue/models"
)

// feat: AuthMiddleware - extract token, find user, and set in context.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		// feat: Expect header format "Bearer user@example.com".
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}

		token := parts[1] // feat: token is user's email in this implementation

		// feat: Find user by email.
		var user models.User
		err := configs.DB.Collection("users").FindOne(context.Background(), bson.M{"email": token}).Decode(&user)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token or user not found"})
			return
		}

		// feat: Set user in context.
		c.Set("user", user)
		c.Next()
	}
}
