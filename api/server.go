package api

import (
	"os"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func Server() {

	router := gin.Default()

	api := router.Group("/api")
	service := Service{}

	reservations := api.Group("/reservations")
	{
		reservations.GET("/:id", service.getReservationByID)
	}

	// serve static files
	router.Use(static.Serve("/", static.LocalFile("./dist", true)))
	router.Run(os.Getenv("HOST"))
}
