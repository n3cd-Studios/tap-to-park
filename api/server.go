package api

import (
	"context"
	"fmt"
	"os"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

var Db *pgx.Conn

func Server() {
	Db, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	defer Db.Close(context.Background())

	router := gin.Default()
	service := Service{}

	api := router.Group("/api")
	api.GET("/test", service.getAlbums) // api: /api/test

	// serve static files
	router.Use(static.Serve("/", static.LocalFile("./dist", true)))

	router.Run(os.Getenv("HOST"))
}
