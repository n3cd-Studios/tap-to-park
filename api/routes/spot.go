package routes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tap-to-park/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

type SpotRoutes struct{}

// GetSpotsNear godoc
//
// @Summary		Get spots near
// @Description	Get a spot near a longitude and latitude
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

	point := database.Coordinates{Longitude: lng, Latitude: lat}
	spots := []database.Spot{}
	query := database.Db.Order(clause.OrderBy{Expression: clause.Expr{SQL: "coords <-> Point (?,?)", Vars: []interface{}{point.Latitude, point.Longitude}, WithoutParentheses: true}}).Limit(10)
	if c.Query("handicap") == "true" {
		query = query.Where("handicap = ?", true)
	}
	result := query.Find(&spots)
	if result.Error != nil {
		c.String(http.StatusNotFound, "Could not load the list of spots.")
		return
	}

	c.IndentedJSON(http.StatusOK, spots)
}

type GetSpotOutput struct {
	database.Spot
	Price float64 `json:"price"`
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
// @Failure		404	{string} string	"Spot was not found"
// @Router		/spots/{id} [get]
func (*SpotRoutes) GetSpot(c *gin.Context) {

	id := c.Param("id")

	spot := database.Spot{}
	result := database.Db.Where("guid = ?", id).First(&spot)
	err := result.Error

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, "Spot was not found")
		return
	}

	c.IndentedJSON(http.StatusOK, GetSpotOutput{
		Spot:  spot,
		Price: spot.GetPrice(),
	})
}

// UpdateSpot godoc
//
// @Summary		Update a spot
// @Description	Update a spot's information such as pricing table, name or longitude and latitude
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

	input := database.Pricing{}
	if err := c.BindJSON(&input); err != nil {
		c.String(http.StatusBadRequest, "Invalid body.")
		return
	}

	id := c.Param("id")
	table, _ := json.Marshal(input)
	if result := database.Db.Model(&database.Spot{}).Where("guid = ?", id).Update("pricing", table); result.Error != nil {
		c.String(http.StatusNotFound, "That spot does not exist.")
		return
	}

	c.String(http.StatusOK, "Successfully updated spot.")
}

type CreateSpotInput struct {
	Name   string               `json:"name" binding:"required"`
	Coords database.Coordinates `json:"coords" binding:"required"`
}

// CreateSpot godoc
//
// @Summary		Create a spot
// @Description	Create a spot at a longitude and latitude
// @Tags		spot
// @Accept		json
// @Produce		json
// @Success		200	{object}	database.Spot
// @Failure		400 {string} string	"Invalid body."
// @Failure		400 {string} string	"Longitude must be between -180 and 180."
// @Failure		400 {string} string	"Latitude must be between -90 and 90."
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

	if input.Coords.Longitude < -180 || input.Coords.Longitude > 180 {
		c.String(http.StatusBadRequest, "Longitude must be between -180 and 180.")
		return
	}
	if input.Coords.Latitude < -90 || input.Coords.Latitude > 90 {
		c.String(http.StatusBadRequest, "Latitude must be between -90 and 90.")
		return
	}

	user := c.MustGet("user").(database.User)
	existingSpot := database.Spot{}
	if err := database.Db.Where("name = ? AND organization_id = ?", input.Name, user.OrganizationID).First(&existingSpot).Error; err == nil {
		c.String(http.StatusConflict, "A spot with this name already exists for the organization.")
		return
	}

	spot := database.Spot{
		Name:           input.Name,
		Coords:         input.Coords,
		OrganizationID: user.OrganizationID,
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
	if result := database.Db.Where("id = ?", spot_id).Where("organization_id = ?", user.OrganizationID).Delete(&spot); result.Error != nil {
		c.String(http.StatusNotFound, "That spot does not exist.")
		return
	}

	c.String(http.StatusOK, "Spot successfully deleted.")
}
