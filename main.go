package main

import (
	"log"
	"tap-to-park/api"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Start our server
	api.Server()
}
