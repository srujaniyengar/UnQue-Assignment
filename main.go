package main

import (
	"UnQue/configs"
	"UnQue/routes"
	"log"
)

func main() {
	configs.ConnectDB()

	router := routes.SetupRoutes()

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
