package routes

import (
	"encoding/json"
	"math"
	"net/http"
	"os"
	"strconv"
	"tap-to-park/database"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
	"gorm.io/gorm/clause"
)

type SpotRoutes struct{}

type ReservationTimes struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type GetSpotsNearOutput struct {
	Guid        string               `json:"guid"`
	Coords      database.Coordinates `json:"coords"`
	Reservation *ReservationTimes    `json:"reservation"`
}

// GetSpotsNear godoc
//
// @Summary		Get spots near
// @Description	Get a spot near a latitude and longitude
// @Tags			spot
// @Accept			json
// @Produce		json
// @Param			lat	query		decimal	true	"Latitude to search near"
// @Param			lng	query		decimal	true	"Longitude to search near"
// @Param			handicap	query		boolean	false	"To filter spots by handicap spots"
// @Success		200	{object}	database.Spot
// @Failure		400	{string} string "Latitude must be a number."
// @Failure		400	{string} string "Longitude must be a number."
// @Failure		404	{string} string "Could not load the list of spots."
// @Router			/spots/near [get]
func (*SpotRoutes) GetSpotsNear(c *gin.Context) {

	latParam := c.Query("lat")
	lngParam := c.Query("lng")

	lat, perr := strconv.ParseFloat(latParam, 64)
	if perr != nil {
		c.String(http.StatusBadRequest, "Latitude must be a number.")
		return
	}

	lng, perr := strconv.ParseFloat(lngParam, 64)
	if perr != nil {
		c.String(http.StatusBadRequest, "Longitude must be a number.")
		return
	}

	spots := []database.Spot{}
	query := database.Db.Model(&database.Spot{}).Order(clause.OrderBy{Expression: clause.Expr{SQL: "coords <-> Point (?,?)", Vars: []interface{}{lat, lng}, WithoutParentheses: true}}).Limit(10)
	if c.Query("handicap") == "true" {
		query = query.Where("handicap = ?", true)
	} else {
		query = query.Where("handicap = ?", false)
	}

	result := query.Find(&spots)
	if result.Error != nil {
		c.String(http.StatusNotFound, "Could not load the list of spots.")
		return
	}

	spotsOutput := []GetSpotsNearOutput{}
	for _, spot := range spots {
		reservation := spot.GetReservation()
		var reservationTimes *ReservationTimes
		if reservation != nil {
			reservationTimes = &ReservationTimes{
				Start: reservation.Start,
				End:   reservation.End,
			}
		}
		spotsOutput = append(spotsOutput, GetSpotsNearOutput{
			Guid:        spot.Guid,
			Coords:      spot.Coords,
			Reservation: reservationTimes,
		})
	}

	c.IndentedJSON(http.StatusOK, spotsOutput)
}

type GetSpotOutput struct {
	database.Spot
	Price       float64               `json:"price"`
	Reservation *database.Reservation `json:"reservation"`
}

// GetSpotsNear godoc
//
// @Summary		Get a spot
// @Description	Get a spot by its ID
// @Tags		spot
// @Accept		json
// @Produce		json
// @Param		id  path		string	true	"The ID of the spot"
// @Success		200	{object}	database.Spot
// @Failure		404	{string} string	"Spot was not found."
// @Router		/spots/{id} [get]
func (*SpotRoutes) GetSpot(c *gin.Context) {

	id := c.Param("id")

	spot := database.Spot{}
	if result := database.Db.Where("guid = ?", id).First(&spot); result.Error != nil {
		c.String(http.StatusNotFound, "Spot was not found.")
		return
	}

	c.IndentedJSON(http.StatusOK, GetSpotOutput{
		Spot:        spot,
		Price:       spot.GetPrice(),
		Reservation: spot.GetReservation(),
	})
}

// GetSpotQRCode godoc
//
// @Summary		Get a spot's QRCode
// @Description	Generates the QRCode that is associated with a spot
// @Tags		spot
// @Accept		json
// @Produce		png
// @Param		id  path		string	true	"The ID of the spot"
// @Success		200	{png} png "The QR Code that was generated"
// @Failure		404	{string} string	"Spot was not found."
// @Failure		500	{string} string	"Failed to generate QR Code."
// @Router		/spots/{id}/qr [get]
func (*SpotRoutes) GetSpotQR(c *gin.Context) {

	id := c.Param("id")

	spot := database.Spot{}
	if result := database.Db.Where("guid = ?", id).First(&spot); result.Error != nil {
		c.String(http.StatusNotFound, "Spot was not found.")
		return
	}

	link := os.Getenv("FRONTEND_HOST") + "/" + spot.Guid
	qr, err := qrcode.Encode(link, qrcode.Medium, 256)

	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to generate QR Code.")
		return
	}

	c.Data(http.StatusOK, "image/png", qr)
}

type UpdateSpotInput struct {
	Table    database.Pricing `json:"table"`
	Name     string           `json:"name"`
	MaxHours uint             `json:"maxHours"`
}

// UpdateSpot godoc
//
// @Summary		Update a spot
// @Description	Update a spot's information such as pricing table, name or latitude and longitude
// @Tags		spot
// @Accept		json
// @Produce		json
// @Param		id  path		string	true	"The ID of the spot"
// @Success		200	{string} string	"Successfully updated spot."
// @Failure		400	{string} string	"Invalid body."
// @Failure		401	{string} string "Unauthorized."
// @Failure		404	{string} string	"That spot does not exist."
// @Router		/spots/{id} [put]
// @Security 	BearerToken
func (*SpotRoutes) UpdateSpot(c *gin.Context) {

	input := UpdateSpotInput{}
	if err := c.BindJSON(&input); err != nil {
		c.String(http.StatusBadRequest, "Invalid body.")
		return
	}

	id := c.Param("id")

	table, _ := json.Marshal(input.Table)
	result := database.Db.Model(&database.Spot{}).Where("guid = ?", id).Update("pricing", table).Update("name", input.Name).Update("max_hours", input.MaxHours)
	if result.Error != nil {
		c.String(http.StatusNotFound, "That spot does not exist.")
		return
	}

	c.String(http.StatusOK, "Successfully updated spot.")
}

type CreateSpotInput struct {
	Name     string               `json:"name" binding:"required"`
	Coords   database.Coordinates `json:"coords" binding:"required"`
	Price    float64              `json:"price" binding:"required"`
	MaxHours uint                 `json:"maxHours" binding:"required"`
	Handicap *bool                `json:"handicap" binding:"required"`
}

// CreateSpot godoc
//
// @Summary		Create a spot
// @Description	Create a spot at a latitude and longitude
// @Tags		spot
// @Accept		json
// @Produce		json
// @Success		200	{object}	database.Spot
// @Failure		400 {string} string	"Invalid body."
// @Failure		400 {string} string	"Latitude must be between -90 and 90."
// @Failure		400 {string} string	"Longitude must be between -180 and 180."
// @Failure		400 {string} string "A spot with this name already exists for the organization."
// @Failure		401	{string} string "Unauthorized."
// @Failure		409 {string} string "A spot with this name already exists for the organization."
// @Router		/spots [post]
// @Security 	BearerToken
func (*SpotRoutes) CreateSpot(c *gin.Context) {

	var input CreateSpotInput
	if err := c.BindJSON(&input); err != nil {
		c.String(http.StatusBadRequest, "Invalid body.")
		return
	}

	if input.Coords.Latitude < -90 || input.Coords.Latitude > 90 {
		c.String(http.StatusBadRequest, "Latitude must be between -90 and 90.")
		return
	}
	if input.Coords.Longitude < -180 || input.Coords.Longitude > 180 {
		c.String(http.StatusBadRequest, "Longitude must be between -180 and 180.")
		return
	}
	if input.Price < 0 {
		c.String(http.StatusBadRequest, "Price must be positive.")
		return
	}
	if math.Round(input.Price*100)/100 != input.Price {
		c.String(http.StatusBadRequest, "Price must have no more than two decimal spaces.")
		return
	}

	user := c.MustGet("user").(database.User)
	existingSpot := database.Spot{}
	if err := database.Db.Where("name = ? AND organization_id = ?", input.Name, user.OrganizationID).First(&existingSpot).Error; err == nil {
		c.String(http.StatusConflict, "A spot with this name already exists for the organization.")
		return
	}

	hours := make([]float64, 24)
	for i := range hours {
		hours[i] = input.Price
	}

	pricing := database.Pricing{
		Monday:    hours,
		Tuesday:   hours,
		Wednesday: hours,
		Thursday:  hours,
		Friday:    hours,
		Saturday:  hours,
		Sunday:    hours,
	}

	spot := database.Spot{
		Name:           input.Name,
		Coords:         input.Coords,
		OrganizationID: user.OrganizationID,
		Pricing:        pricing,
		Handicap:       *input.Handicap,
		MaxHours:       input.MaxHours,
	}
	if result := database.Db.Create(&spot); result.Error != nil {
		c.String(http.StatusBadRequest, "Failed to create a spot.")
		return
	}

	c.IndentedJSON(http.StatusOK, spot)
}

// DeleteSpot godoc
//
// @Summary		Delete a spot
// @Description	Delete a spot by its ID
// @Tags		spot
// @Accept		json
// @Produce		json
// @Param		id  path		string	true	"The ID of the spot"
// @Success		200	{string} string	"Spot successfully deleted."
// @Failure		401	{string} string "Unauthorized."
// @Failure		404 {string} string	"That spot does not exist."
// @Router		/spots/{id} [delete]
// @Security 	BearerToken
func (*SpotRoutes) DeleteSpot(c *gin.Context) {

	user := c.MustGet("user").(database.User)
	spot_id := c.Param("id")
	spot := database.Spot{}

	if result := database.Db.Where("guid = ?", spot_id).Where("organization_id = ?", user.OrganizationID).Delete(&spot); result.Error != nil {
		c.String(http.StatusInternalServerError, "Failed to delete the spot.")
		return
	}

	c.IndentedJSON(http.StatusOK, "Spot successfully deleted.")
}
