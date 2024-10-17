package routes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"tap-to-park/database"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

type SpotRoutes struct{}

// GetSpotsNear godoc
// @Summary      Get the spots near a longitude and latitude, with optional handicap filter
// @Produce      json
// @Param        lat    query     number  true  "latitude to search by"
// @Param        lng    query     number  true  "longitude to search by"
// @Param        handicap query     boolean false "filter spots by handicap accessibility"
// @Success      200  {object}  database.Spot
// @Failure      404  {string}  "Could not load the list of spots"
// @Router       /spots/near [get]
func (*SpotRoutes) GetSpotsNear(c *gin.Context) {

	latParam := c.Query("lat")
	lngParam := c.Query("lng")
	handicapParam := c.Query("handicap")

	lat, perr := strconv.ParseFloat(latParam, 64)
	if perr != nil {
		c.String(http.StatusNotFound, "Latitude must be a float.")
		return
	}

	lng, perr := strconv.ParseFloat(lngParam, 64)
	if perr != nil {
		c.String(http.StatusNotFound, "Longitude must be a float.")
		return
	}

	var filterByHandicap, applyHandicapFilter bool
	if handicapParam != "" {
		filterByHandicap, perr = strconv.ParseBool(handicapParam)
		if perr != nil {
			c.String(http.StatusNotFound, "Handicap must be a boolean.")
			return
		}
		applyHandicapFilter = true
	}

	point := database.Coordinates{Longitude: lng, Latitude: lat}
	spots := []database.Spot{}
	query := database.Db.Order(clause.OrderBy{Expression: clause.Expr{SQL: "coords <-> Point (?,?)", Vars: []interface{}{point.Latitude, point.Longitude}, WithoutParentheses: true}}).Limit(10)
	if applyHandicapFilter {
		query = query.Where("handicap = ?", filterByHandicap)
	}
	result := query.Find(&spots)
	err := result.Error

	if err != nil {
		c.String(http.StatusNotFound, "Could not load the list of spots.")
		return
	}

	c.IndentedJSON(http.StatusAccepted, spots)
}

// GetSpot godoc
// @Summary      Get the spots near a longitude and latitude
// @Produce      json
// @Param        id    path     uuid  true  "Guid of the spot"
// @Success      200  {object}  database.Spot
// @Failure      404  {string}  "Spot was not found"
// @Router       /spots/{id} [get]
func (*SpotRoutes) GetSpot(c *gin.Context) {

	id := c.Param("id")

	spot := database.Spot{}
	result := database.Db.Where("guid = ?", id).First(&spot)
	err := result.Error

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, "Spot was not found")
		return
	}

	now := time.Now()
	weekday := strings.ToLower(now.Weekday().String())
	hour := strconv.FormatInt(int64(now.Hour()), 10)

	// WHY DOES THIS WORK
	var price float64
	if result := database.Db.Raw("select ((pricing->?)::JSON->>CAST(? AS INT))::DECIMAL as price from spots where id=?", weekday, hour, spot.ID).Scan(&price); result.Error != nil {
		return
	}

	println(price)

	c.IndentedJSON(http.StatusOK, spot)
}

// UpdateSpot godoc
// @Summary      Update a spot by its ID
// @Success      200  {string}  "Successfully deleted spot"

// @Failure      400  {string}  "Invalid body"
// @Failure      401  {string}  "Invalid token"
// @Router       /spots/{id} [put]
func (*SpotRoutes) UpdateSpot(c *gin.Context) {

	guid := c.MustGet("guid")

	user := database.User{}
	if result := database.Db.Where("guid = ?", guid).First(&user); result.Error != nil {
		c.String(http.StatusBadRequest, "You literally don't exist")
		return
	}

	input := database.Pricing{}
	if err := c.BindJSON(&input); err != nil {
		c.String(http.StatusBadRequest, "Invalid body")
		return
	}

	id := c.Param("id")
	table, _ := json.Marshal(input)
	if result := database.Db.Model(&database.Spot{}).Where("guid = ?", id).Update("pricing", table); result.Error != nil {
		c.String(http.StatusNotFound, "That spot does not exist")
		return
	}

	c.String(http.StatusOK, "Successfully updated spot")
}

type CreateSpotInput struct {
	Name   string               `json:"name" binding:"required"`
	Coords database.Coordinates `json:"coords" binding:"required"`
}

// CreateSpot godoc
// @Summary      Create a spot at a longitude and latitude
// @Produce      json
// @Accept		 json
// @Success      200  {object}  database.Spot
// @Failure      400  {string}  "Invalid body"
// @Failure      401  {string}  "Invalid token"
// @Router       /spots [post]
func (*SpotRoutes) CreateSpot(c *gin.Context) {

	guid := c.MustGet("guid")

	user := database.User{}
	if result := database.Db.Where("guid = ?", guid).First(&user); result.Error != nil {
		c.String(http.StatusNotFound, "For some reason, you don't exist!")
		return
	}

	var input CreateSpotInput
	if err := c.BindJSON(&input); err != nil {
		c.String(http.StatusBadRequest, "Invalid body")
		return
	}

	if input.Coords.Longitude < -180 || input.Coords.Longitude > 180 {
		c.String(http.StatusBadRequest, "Longitude must be between -180 and 180")
		return
	}
	if input.Coords.Latitude < -90 || input.Coords.Latitude > 90 {
		c.String(http.StatusBadRequest, "Latitude must be between -90 and 90")
		return
	}

	existingSpot := database.Spot{}
	if err := database.Db.Where("name = ? AND organization_id = ?", input.Name, user.OrganizationID).First(&existingSpot).Error; err == nil {
		c.String(http.StatusBadRequest, "A spot with this name already exists for the organization")
		return
	}

	spot := database.Spot{
		Name:           input.Name,
		Coords:         input.Coords,
		OrganizationID: user.OrganizationID,
	}
	if result := database.Db.Create(&spot); result.Error != nil {
		c.String(http.StatusBadRequest, "Failed to create a spot")
		return
	}

	c.IndentedJSON(http.StatusOK, spot)
}

// DeleteSpot godoc
// @Summary      Delete a spot by its ID
// @Success      200  {string}  "Successfully deleted spot"
// @Failure      400  {string}  "Invalid body"
// @Failure      401  {string}  "Invalid token"
// @Router       /spots/{id} [delete]
func (*SpotRoutes) DeleteSpot(c *gin.Context) {

	guid := c.MustGet("guid")

	user := database.User{}
	if result := database.Db.Where("guid = ?", guid).First(&user); result.Error != nil {
		c.String(http.StatusBadRequest, "You literally don't exist")
		return
	}

	spot_id := c.Param("id")

	spot := database.Spot{}
	if result := database.Db.Where("id = ?", spot_id).Where("organization_id = ?", user.OrganizationID).Delete(&spot); result.Error != nil {
		c.String(http.StatusNotFound, "That spot does not exist")
		return
	}

	c.String(http.StatusOK, "Spot successfully deleted")
}
