package routes

import (
	"net/http"
	"strconv"
	"tap-to-park/database"

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

// GetSpotByID godoc
// @Summary      Get the spots near a longitude and latitude
// @Produce      json
// @Param        id    query     uuid  true  "Guid of the spot"
// @Success      200  {object}  database.Spot
// @Failure      404  {string}  "Spot was not found"
// @Router       /spots/info [get]
func (*SpotRoutes) GetSpotByID(c *gin.Context) {

	guid := c.Query("guid")

	spot := database.Spot{}
	result := database.Db.Where("guid = ?", guid).First(&spot)
	err := result.Error

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, "Spot was not found")
		return
	}

	c.IndentedJSON(http.StatusAccepted, spot)
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
// @Router       /spots/create [post]
func (*SpotRoutes) CreateSpot(c *gin.Context) {

	uuid := c.MustGet("uuid")

	var input CreateSpotInput
	if err := c.BindJSON(&input); err != nil {
		c.String(http.StatusBadRequest, "Invalid body")
		return
	}

	user := database.User{}
	if result := database.Db.Where("uuid = ?", uuid).First(&user); result.Error != nil {
		c.String(http.StatusBadRequest, "You literally don't exist")
		return
	}

	spot := database.Spot{
		Name:           input.Name,
		Coords:         input.Coords,
		OrganizationID: 1,
	}
	if result := database.Db.Create(&spot); result.Error != nil {
		c.String(http.StatusBadRequest, "Failed to create a spot")
		return
	}

	c.IndentedJSON(http.StatusAccepted, spot)
}

type DeleteSpotInput struct {
	SpotID uint64 `json:"spot_id" bindings:"required"`
}

// DeleteSpot godoc
// @Summary      Delete a spot by it's ID
// @Success      200  {string}  "Successfully deleted spot"
// @Failure      400  {string}  "Invalid body"
// @Failure      401  {string}  "Invalid token"
// @Router       /spots/delete [delete]
func (*SpotRoutes) DeleteSpot(c *gin.Context) {

	uuid := c.MustGet("uuid")

	var input DeleteSpotInput
	if err := c.BindJSON(&input); err != nil {
		c.String(http.StatusBadRequest, "Invalid body")
		return
	}

	// TODO: Make this more efficient

	user := database.User{}
	if result := database.Db.Where("uuid = ?", uuid).First(&user); result.Error != nil {
		c.String(http.StatusBadRequest, "You literally don't exist")
		return
	}

	spot := database.Spot{}
	if result := database.Db.Where("id = ?", input.SpotID).Where("organization_id = ?", user.OrganizationID).Delete(&spot); result.Error != nil {
		c.String(http.StatusNotFound, "That spot does not exist")
		return
	}

	c.String(http.StatusAccepted, "Spot successfully deleted")
}
