package api

import (
	"os"
	"tap-to-park/api/services"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func Server() {

	router := gin.Default()

	api := router.Group("/api")

	reservations := api.Group("/reservations")
	{
		service := services.ReservationService{}
		reservations.GET("/:id", service.getReservationByID)
	}

	// serve static files
	router.Use(static.Serve("/", static.LocalFile("./dist", true)))

	router.Run(os.Getenv("HOST"))
}
