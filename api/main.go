package main

import (
	"log"
	"os"
	"tap-to-park/database"
	"tap-to-park/routes"

	_ "tap-to-park/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Tap-To-Park API
// @version         1.0
// @description     This is the API for interacting with internal Tap-To-Park services
// @termsOfService  http://n3cd.io/terms/

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {

	// Load .env file
	err := godotenv.Load("../.env")
	if err != nil {
		log.Panic("Error loading .env file")
	}

	// Connect to the database
	database.Connect()

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowHeaders:    []string{"Authentication"},
	}))

	// API
	api := router.Group("/api")

	// Reservation routes
	reservations := api.Group("/reservations")
	{
		routing := routes.ReservationRoutes{}
		reservations.POST("/", routing.CreateReservation)
		reservations.GET("/:id", routing.GetReservationByID)
	}

	// Spot routes
	spots := api.Group("/spots")
	{
		routing := routes.SpotRoutes{}
		spots.GET("/near", routing.GetSpotsNear)
	}

	// Auth routes
	auth := api.Group("/auth")
	{
		routing := routes.AuthRoutes{}
		auth.POST("/login", routing.Login)
		auth.POST("/register", routing.Register)
		auth.GET("/info", routes.AuthMiddleware(), routing.Info)
	}

	// Auth routes
	admin := api.Group("/admin", routes.AuthMiddleware())
	{
		routing := routes.AdminRoutes{}
		admin.GET("/organization", routing.GetOrganization)
	}

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	router.Run(os.Getenv("BACKEND_HOST"))
}
