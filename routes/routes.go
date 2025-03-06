package routes

import (
	"UnQue/control"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()

	router.POST("/login", control.Login)

	router.POST("/availability", control.SetAvailability)

	router.GET("/availability", control.GetAvailability)

	router.POST("/appointments", control.BookAppointment)

	router.DELETE("/appointments/:id", control.CancelAppointment)

	router.GET("/appointments", control.GetAppointments)

	return router
}
