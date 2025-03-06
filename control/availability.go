package control

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"UnQue/configs"
	"UnQue/models"
)

func SetAvailability(c *gin.Context) {

	usr_intf, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	proff, ok := usr_intf.(models.User)
	if !ok || proff.Role != "proff" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only professors can set availability"})
		return
	}

	var input struct {
		Slots []string `json:"slots"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	contxt, cnl := context.WithTimeout(context.Background(), 10*time.Second)
	defer cnl()

	slot_collection := configs.DB.Collection("timeslots")
	var createdSlots []models.TimeSlot

	for _, slot := range input.Slots {
		timeslot := models.TimeSlot{
			Professor: proff.ID,
			Slot:      slot,
			Booked:    false,
		}
		result, err := slot_collection.InsertOne(contxt, timeslot)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set availability"})
			return
		}
		timeslot.ID = result.InsertedID.(primitive.ObjectID)
		createdSlots = append(createdSlots, timeslot)
	}

	c.JSON(http.StatusOK, createdSlots)
}

func GetAvailability(c *gin.Context) {
	professorIDHex := c.Query("professor_id")
	if professorIDHex == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "professor_id is required"})
		return
	}

	professorID, err := primitive.ObjectIDFromHex(professorIDHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid proff ID"})
		return
	}

	contxt, cnl := context.WithTimeout(context.Background(), 10*time.Second)
	defer cnl()

	slot_collection := configs.DB.Collection("timeslots")
	filter := bson.M{"proff": professorID, "booked": false}
	cursor, err := slot_collection.Find(contxt, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch availability"})
		return
	}

	var slots []models.TimeSlot
	if err = cursor.All(contxt, &slots); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding slots"})
		return
	}

	c.JSON(http.StatusOK, slots)
}
