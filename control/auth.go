package control

import (
	"context"
	"log" // <-- Added for debug logs
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"

	"UnQue/configs"
	"UnQue/models"
)

// feat: Login - handle POST /login requests.
func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// feat: Bind JSON input for login.
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("fix: Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// feat: Log login attempt.
	log.Printf("feat: Login attempt: email=%s, password=%s", input.Email, input.Password)

	// feat: Create context with timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userCollection := configs.DB.Collection("users")
	var user models.User

	// feat: Query user by email and password.
	err := userCollection.FindOne(ctx, bson.M{
		"email":    input.Email,
		"password": input.Password,
	}).Decode(&user)
	log.Printf("feat: Using database: %s, collection: users", configs.DB.Name())

	// fix: Handle error if user is not found.
	if err != nil {
		log.Printf("fix: FindOne error: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// feat: Return token as email and user info on successful login.
	log.Printf("feat: Login success for user: %v", user.Email)
	c.JSON(http.StatusOK, gin.H{
		"token": user.Email,
		"user":  user,
	})
}
