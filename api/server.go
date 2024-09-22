package main

import (
	"log"
	"os"
	"tap-to-park/database"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	// Load .env file
	err := godotenv.Load("../.env")
	if err != nil {
		log.Panic("Error loading .env file")
	}

	// Connect to the database
	database.Connect()

	router := gin.Default()

	api := router.Group("/api")
	service := Service{}

	reservations := api.Group("/reservations")
	{
		reservations.POST("/", service.postReservation)
		reservations.GET("/:id", service.getReservationByID)
	}

	router.Run(os.Getenv("BACKEND_HOST"))
}
