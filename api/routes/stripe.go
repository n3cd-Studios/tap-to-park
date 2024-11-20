package routes

import (
	"net/http"
	"os"
	"strconv"
	"tap-to-park/database"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v80"
	"github.com/stripe/stripe-go/v80/checkout/session"
)

type StripeRoutes struct{}

// SuccessfulPurchaseSpot godoc
//
// @Summary		Success callback for Stripe
// @Description	Create a Reservation from a Stripe Session ID
// @Tags		spot,reservation,stripe
// @Accept		json
// @Produce		json
// @Param		id	path		string true	"The ID of the Spot"
// @Param		session_id	query		string	true	"The Session ID passed from Stripe"
// @Success		301	{string} string "This will redirect you to a page on the frontend."
// @Failure		400	{string} string "That spot ID is invalid."
// @Failure		400	{string} string "Invalid Stripe session."
// @Failure		400	{string} string "Something went wrong (did you resubmit the request?)"
// @Failure		500	{string} string "Malformed metadata."
// @Router		/stripe/{id}/success [get]
func (*StripeRoutes) SuccessfulPurchaseSpot(c *gin.Context) {

	spot_id := c.Param("id")
	spot := database.Spot{}
	if result := database.Db.Where("guid = ?", spot_id).First(&spot); result.Error != nil {
		c.String(http.StatusBadRequest, "That spot ID is invalid.")
		return
	}

	if spot.GetReservation() != nil {
		c.String(http.StatusConflict, "This spot has already been reserved.")
		return
	}

	session_id := c.Query("session_id")
	if session_id == "" {
		c.String(http.StatusBadRequest, "Invalid Stripe session.")
		return
	}

	sess, err := session.Get(session_id, nil)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid Stripe session.")
		return
	}

	hours, err := strconv.ParseFloat(sess.Metadata["hours"], 64)
	if err != nil {
		c.String(http.StatusInternalServerError, "Malformed metadata.")
		return
	}

	now := time.Now()
	reservation := database.Reservation{
		Start:               now,
		End:                 now.Add(time.Duration(hours*60) * time.Minute),
		Email:               sess.CustomerDetails.Email,
		Price:               float64(sess.AmountTotal),
		SpotID:              spot.ID,
		StripeTransactionID: sess.ID,
	}

	// We can possibly assign this transaction to a user.
	user_id, err := strconv.ParseUint(sess.Metadata["user_id"], 10, 32)
	if err == nil {
		reservation.UserID = uint(user_id)
	}

	if result := database.Db.Create(&reservation); result.Error != nil {
		c.String(http.StatusBadRequest, "Something went wrong (did you resubmit the request?)")
		return
	}

	c.Redirect(http.StatusMovedPermanently, os.Getenv("FRONTEND_HOST")+"/"+spot.Guid+"/success")
}

// CancelPurchaseSpot godoc
//
// @Summary		Cancel callback for Stripe
// @Description	This is just a dummy route, it redirects to the frontend
// @Tags		spot,reservation,stripe
// @Accept		json
// @Produce		json
// @Param		id	path		string true	"The ID of the Spot"
// @Success		301	{string} string "This will redirect you to a page on the frontend."
// @Failure		400	{string} string "That spot ID is invalid."
// @Router		/stripe/{id}/cancel [get]
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
	User  string `json:"user"`
}

// PurchaseSpot godoc
//
// @Summary		Create a checkout session
// @Description	Create a Stripe checkout session and forward the user to Stripe
// @Tags		spot,reservation,stripe
// @Accept		json
// @Produce		json
// @Param		id	path		string true	"The ID of the Spot"
// @Param		session_id	query		string	true	"The Session ID passed from Stripe"
// @Success		200	{string} string "The Stripe checkout URL"
// @Failure		400	{string} string "That spot ID is invalid."
// @Failure		400	{string} string "Invalid body."
// @Failure		400 {string} string "You can't purchase a spot for this amount of time."
// @Failure		400 {string} string "Reservation exceeds spot maximum reservation time."
// @Failure		400 {string} string "Reservation cost must be at least 50¢."
// @Failure		409 {string} string "This spot has already been reserved."
// @Failure		500	{string} string "Failed to create stripe session."
// @Router		/stripe/{id} [post]
func (*StripeRoutes) PurchaseSpot(c *gin.Context) {

	var input PurchaseSpotInput
	if err := c.BindJSON(&input); err != nil {
		c.String(http.StatusBadRequest, "Invalid body.")
		return
	}

	spot_id := c.Param("id")
	spot := database.Spot{}
	if result := database.Db.Where("guid = ?", spot_id).First(&spot); result.Error != nil {
		c.String(http.StatusBadRequest, "That spot ID is invalid")
		return
	}

	if spot.GetReservation() != nil {
		c.String(http.StatusConflict, "This spot has already been reserved.")
		return
	}

	hours := input.End.Sub(input.Start).Hours()
	if hours < 0 {
		c.String(http.StatusBadRequest, "You can't purchase a spot for this amount of time.")
		return
	}

	if hours > float64(spot.MaxHours) {
		c.String(http.StatusBadRequest, "Reservation exceeds spot maximum reservation time.")
		return
	}

	var purchasePrice = int64(hours * spot.GetPrice() * 100)
	if purchasePrice < 50 {
		c.String(http.StatusBadRequest, "Reservation cost must be at least 50¢.")
		return
	}

	// TODO: we could check if this is valid and pass an auth token, but why would someone want to randomly buy
	// a spot for another user??
	user := database.User{}
	metadata := map[string]string{"hours": strconv.FormatFloat(hours, 'E', -1, 64)}
	if result := database.Db.Where("guid = ?", input.User).First(&user); result.Error == nil {
		metadata["user_id"] = strconv.FormatUint(uint64(user.ID), 64)
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
					UnitAmount: stripe.Int64(purchasePrice),
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(domain + "/api/stripe/" + spot.Guid + "/success?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String(domain + "/api/stripe/" + spot.Guid + "/cancel"),
		Metadata:   metadata,
	}

	sess, err := session.New(params)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to create Stripe session.")
		return
	}

	c.IndentedJSON(http.StatusOK, sess)
}
