package api

import (
	"os"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func Server() {
	router := gin.Default()
	service := Service{}

	api := router.Group("/api")
	api.GET("/test", service.getAlbums) // api: /api/test

	// serve static files
	router.Use(static.Serve("/", static.LocalFile("./dist", true)))

	router.Run(os.Getenv("HOST"))
}
