package control

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"

	"UnQue/configs"
	"UnQue/models"
)

// Login handles POST /login requests.
func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	ctx, cnl := context.WithTimeout(context.Background(), 10*time.Second)
	defer cnl()

	user_collection := configs.DB.Collection("users")
	var user models.User
	err := user_collection.FindOne(ctx, bson.M{"username": input.Username, "password": input.Password}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// In this example, the token is simply the username.
	c.JSON(http.StatusOK, gin.H{
		"token": user.Username,
		"user":  user,
	})
}
