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

// BookAppointment handles POST /appointments.
func BookAppointment(c *gin.Context) {
	usr_intf, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	student, ok := usr_intf.(models.User)
	if !ok || student.Role != "student" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only students can book appointments"})
		return
	}

	var input struct {
		ProfessorID string `json:"professor_id"`
		Slot        string `json:"slot"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	proff_obj_id, err := primitive.ObjectIDFromHex(input.ProfessorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid professor ID"})
		return
	}

	get_slots, cnl := context.WithTimeout(context.Background(), 10*time.Second)
	defer cnl()

	slot_collection := configs.DB.Collection("timeslots")
	filer := bson.M{"professor": proff_obj_id, "slot": input.Slot, "booked": false}
	update := bson.M{"$set": bson.M{"booked": true}}

	var time_slot models.TimeSlot
	err = slot_collection.FindOneAndUpdate(get_slots, filer, update).Decode(&time_slot)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Requested slot not available"})
		return
	}

	appointment_collection := configs.DB.Collection("appointments")
	appointment := models.Appointment{
		Student:      student.ID,
		Professor:    proff_obj_id,
		Availability: time_slot.ID,
	}

	result, err := appointment_collection.InsertOne(get_slots, appointment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create appointment"})
		return
	}

	appointment.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusOK, appointment)
}

// CancelAppointment handles DELETE /appointments/:id.
func CancelAppointment(c *gin.Context) {
	usr_intf, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	professor, ok := usr_intf.(models.User)
	if !ok || professor.Role != "professor" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only professors can cancel appointments"})
		return
	}

	appointmentIDHex := c.Param("id")
	appointmentID, err := primitive.ObjectIDFromHex(appointmentIDHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid appointment ID"})
		return
	}

	contxt, cnl := context.WithTimeout(context.Background(), 10*time.Second)
	defer cnl()

	appointment_collection := configs.DB.Collection("appointments")
	var appointment models.Appointment
	err = appointment_collection.FindOne(contxt, bson.M{"_id": appointmentID, "professor": professor.ID}).Decode(&appointment)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
		return
	}

	_, err = appointment_collection.DeleteOne(contxt, bson.M{"_id": appointmentID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not cancel appointment"})
		return
	}

	timeslot_collection := configs.DB.Collection("timeslots")
	_, err = timeslot_collection.UpdateOne(contxt, bson.M{"_id": appointment.Availability}, bson.M{"$set": bson.M{"booked": false}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Appointment cancelled but failed to update timeslot"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Appointment cancelled"})
}

// GetAppointments handles GET /appointments.
func GetAppointments(c *gin.Context) {
	usr_intf, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	student, ok := usr_intf.(models.User)
	if !ok || student.Role != "student" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only students can view appointments"})
		return
	}

	contxt, cnl := context.WithTimeout(context.Background(), 10*time.Second)
	defer cnl()

	appointment_collection := configs.DB.Collection("appointments")
	cursor, err := appointment_collection.Find(contxt, bson.M{"student": student.ID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching appointments"})
		return
	}

	var appointments []models.Appointment
	if err = cursor.All(contxt, &appointments); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding appointments"})
		return
	}

	c.JSON(http.StatusOK, appointments)
}
