package main

import (
	"log"
	"tap-to-park/api"
	"tap-to-park/api/database"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to the database
	database.Connect()

	// Start our server
	api.Server()
}
