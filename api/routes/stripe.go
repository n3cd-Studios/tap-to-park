package routes

import (
	"net/http"
	"os"
	"tap-to-park/database"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v80"
	"github.com/stripe/stripe-go/v80/checkout/session"
)

type StripeRoutes struct{}

func (*StripeRoutes) SuccessfulPurchaseSpot(c *gin.Context) {

	spot_id := c.Param("id")
	spot := database.Spot{}
	if result := database.Db.Where("guid = ?", spot_id).First(&spot); result.Error != nil {
		c.String(http.StatusBadRequest, "That spot ID is invalid")
		return
	}

	session_id := c.Query("session_id")
	if session_id == "" {
		c.String(http.StatusBadRequest, "Invalid Stripe session")
		return
	}

	sess, err := session.Get(session_id, nil)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid Stripe session")
		return
	}

	reservation := database.Reservation{
		Start:               time.Now(),
		End:                 time.Now(),
		Email:               sess.CustomerDetails.Email,
		Price:               float64(sess.AmountTotal),
		SpotID:              spot.ID,
		StripeTransactionID: sess.ID,
	}
	if result := database.Db.Create(&reservation); result.Error != nil {
		c.String(http.StatusBadRequest, "Something went wrong (did you resubmit the request?)")
		return
	}

	c.Redirect(http.StatusMovedPermanently, os.Getenv("FRONTEND_HOST")+"/"+spot.Guid+"/success")
}

func (*StripeRoutes) CancelPurchaseSpot(c *gin.Context) {

	spot_id := c.Param("id")
	spot := database.Spot{}
	if result := database.Db.Where("guid = ?", spot_id).First(&spot); result.Error != nil {
		c.String(http.StatusBadRequest, "That spot ID is invalid")
		return
	}

	c.Redirect(http.StatusMovedPermanently, os.Getenv("FRONTEND_HOST")+"/"+spot.Guid)
}

type PurchaseSpotInput struct {
	Start time.Time `json:"start" bindings:"required"`
	End   time.Time `json:"end" bindings:"required"`
}

// PurchaseSpot godoc
// @Summary      Purchase a spot with a certain price
// @Produce      json
// @Param        id    path     uuid  true  "Guid of the spot"
// @Success      200  {string}  "URL of the stripe checkout"
// @Failure      400  {string}  "Invalid body"
// @Failure      404  {string}  "Spot was not found"
// @Router       /spots/{id}/purchase [post]
func (*StripeRoutes) PurchaseSpot(c *gin.Context) {

	var input PurchaseSpotInput
	if err := c.BindJSON(&input); err != nil {
		c.String(http.StatusBadRequest, "Invalid body")
		return
	}

	spot_id := c.Param("id")
	spot := database.Spot{}
	if result := database.Db.Where("guid = ?", spot_id).First(&spot); result.Error != nil {
		c.String(http.StatusBadRequest, "That spot ID is invalid")
		return
	}

	domain := "http://" + os.Getenv("BACKEND_HOST")
	params := &stripe.CheckoutSessionParams{
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("usd"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name:        stripe.String("Parking"),
						Description: stripe.String("Parking at " + spot.Name),
					},
					UnitAmount: stripe.Int64(int64(input.End.Sub(input.Start).Hours() * spot.GetPrice() * 100)),
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(domain + "/api/stripe/" + spot.Guid + "/success?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String(domain + "/api/stripe/" + spot.Guid + "/cancel"),
	}

	sess, err := session.New(params)
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to create stripe session.")
		return
	}

	c.IndentedJSON(http.StatusOK, sess)
}
