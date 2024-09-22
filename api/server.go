package main

import (
	"log"
	"os"
	"tap-to-park/database"
	"tap-to-park/routes"

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

	reservations := api.Group("/reservations")
	{
		routes := routes.ReservationRoutes{}
		reservations.POST("/", routes.CreateReservation)
		reservations.GET("/:id", routes.GetReservationByID)
	}

	spots := api.Group("/spots")
	{
		routes := routes.SpotRoutes{}
		spots.GET("/near", routes.GetSpotsNear)
	}

	router.Run(os.Getenv("BACKEND_HOST"))
}
