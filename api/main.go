package main

import (
	"log"
	"os"
	"tap-to-park/database"
	"tap-to-park/routes"
	"tap-to-park/routes/auth"

	_ "tap-to-park/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v80"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Tap-To-Park API
// @version         1.0
// @description     This is the API for interacting with internal Tap-To-Park services
// @termsOfService  http://n3cd.io/terms/
//
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
//
// @host      localhost:8080
// @BasePath  /api/
//
// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {

	// Load .env file
	if !IsProduction() {
		err := godotenv.Load("../.env")
		if err != nil {
			log.Panic("Error loading .env file")
		}
	}

	stripe.Key = os.Getenv("STRIPE_API_KEY")

	// Connect to the database
	database.Connect()

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"PUT", "POST", "OPTIONS", "GET", "DELETE"},
		AllowHeaders:    []string{"Authentication"},
	}))

	// API
	api := router.Group("/api")

	// Reservation routes
	reservations := api.Group("/reservations")
	{
		routing := routes.ReservationRoutes{}
		reservations.GET("/:id", routing.GetReservation)
	}

	// Spot routes
	spots := api.Group("/spots")
	{
		routing := routes.SpotRoutes{}
		spots.GET("/near", routing.GetSpotsNear)
		spots.GET("/:id", routing.GetSpot)
		spots.GET("/:id/qr", routing.GetSpotQR)
		spots.POST("", auth.AuthMiddleware(database.ADMIN), routing.CreateSpot)
		spots.PUT("/:id", auth.AuthMiddleware(database.ADMIN), routing.UpdateSpot)
		spots.DELETE("/:id", auth.AuthMiddleware(database.ADMIN), routing.DeleteSpot)
	}

	// Stripe routes
	stripe := api.Group("/stripe")
	{
		routing := routes.StripeRoutes{}
		stripe.POST("/:id", routing.PurchaseSpot)
		stripe.GET("/:id/success", routing.SuccessfulPurchaseSpot)
		stripe.GET("/:id/cancel", routing.CancelPurchaseSpot)
	}

	// Auth routes
	authr := api.Group("/auth")
	{
		routing := auth.AuthRoutes{}

		// regular login
		authr.POST("/login", routing.Login)
		authr.POST("/register", routing.Register)

		// oauth login
		authr.GET("/:type", routing.OAuthInitialize)
		authr.GET("/:type/callback", routing.OAuthCallback)

		authr.GET("/info", auth.AuthMiddleware(database.ADMIN, database.USER), routing.GetInfo)
		authr.GET("/sessions", auth.AuthMiddleware(database.ADMIN, database.USER), routing.GetSessions)
		authr.DELETE("/sessions/:id", auth.AuthMiddleware(database.ADMIN, database.USER), routing.RevokeSession)
	}

	// Organization routes
	organization := api.Group("/organization", auth.AuthMiddleware(database.ADMIN))
	{
		routing := routes.OrganizationRoutes{}
		organization.GET("/me", routing.GetOrganization)
		organization.GET("/spots", routing.GetSpots)
		organization.GET("/invites", routing.GetInvites)
		organization.POST("/invites", routing.CreateInvite)
	}

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	router.Run(os.Getenv("BACKEND_HOST"))
}

func IsProduction() bool {
	return os.Getenv("GO_ENV") == "production"
}
