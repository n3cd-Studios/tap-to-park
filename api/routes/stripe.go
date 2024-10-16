package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v80"
	"github.com/stripe/stripe-go/v80/checkout/session"
)

type StripeRoutes struct{}

func (*StripeRoutes) SuccessfulPurchaseSpot(c *gin.Context) {
	c.String(http.StatusOK, "Good job")
}

func (*StripeRoutes) CancelPurchaseSpot(c *gin.Context) {
	c.String(http.StatusOK, "Bad job")
}

type PurchaseSpotInput struct {
	Price float64 `json:"price" bindings:"required"`
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

	id := c.Param("id")

	domain := "https://localhost:8080"
	params := &stripe.CheckoutSessionParams{
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				// Provide the exact Price ID (for example, pr_1234) of the product you want to sell
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("usd"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name:        stripe.String("Parking"),
						Description: stripe.String("Parking for " + id),
					},
					UnitAmountDecimal: stripe.Float64(input.Price),
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(domain + "/stripe/" + id + "/success"),
		CancelURL:  stripe.String(domain + "/stripe/" + id + "/cancel"),
	}

	sess, err := session.New(params)
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to create stripe session.")
		return
	}

	c.IndentedJSON(http.StatusOK, sess)
}
