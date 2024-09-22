package main

import (
	"log"
	"os"
	"tap-to-park/database"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func Server() {

	// Load .env file
	err := godotenv.Load()
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

	// serve static files
	router.Use(static.Serve("/", static.LocalFile("./dist", true)))
	router.Run(os.Getenv("HOST"))
}
