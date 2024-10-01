package routes

import (
	"net/http"
	"tap-to-park/database"

	"github.com/gin-gonic/gin"
)

type OrganizationRoutes struct{}

// GetOrganization godoc
// @Summary      Get all of the organizations associated with an admin
// @Produce      json
// @Success      200  {array} []database.Organization
// @Failure      400  {string}  "Unauthorized"
// @Router       /admin/organization [get]
func (*OrganizationRoutes) GetOrganization(c *gin.Context) {

	uuid := c.MustGet("guid")

	user := database.User{}
	if result := database.Db.Where("guid = ?", uuid).First(&user); result.Error != nil {
		c.String(http.StatusNotFound, "For some reason, you don't exist!")
		return
	}

	organization := database.Organization{}
	result := database.Db.Model(&database.Organization{}).Preload("Spots").Where("id = ?", user.OrganizationID).First(&organization)
	if result.Error != nil {
		c.String(http.StatusNotFound, "Couldn't find the organization associated with you")
		return
	}

	c.IndentedJSON(http.StatusOK, organization)
}

// GetSpotData godoc
// @Summary      Get all of the spots data associated with an organization
// @Produce      json
// @Success      200  {array} []database.Spot
// @Failure      400  {string}  "Unauthorized"
// @Router       /admin/organization/data [get]
func (*OrganizationRoutes) GetSpotData(c *gin.Context) {

	uuid := c.MustGet("guid")

	user := database.User{}
	if result := database.Db.Where("guid = ?", uuid).First(&user); result.Error != nil {
		c.String(http.StatusNotFound, "For some reason, you don't exist!")
		return
	}

	organization := database.Organization{}
	orgResult := database.Db.Model(&database.Organization{}).Preload("Spots").Where("id = ?", user.OrganizationID).First(&organization)
	if orgResult.Error != nil {
		c.String(http.StatusNotFound, "Couldn't find the organization associated with you")
		return
	}

	// Gather just the spot IDs from the organization struct
	var spotIDs []uint
	for _, spot := range organization.Spots {
		spotIDs = append(spotIDs, spot.ID)
	}

	var spots []database.Spot
	spotResult := database.Db.Where("id in ?", spotIDs).Find(&spots)
	if spotResult.Error != nil || len(spots) == 0 {
		c.String(http.StatusNotFound, "Couldn't find the spots associated with the organization")
		return
	}

	var reservations []database.Reservation
	reservationResult := database.Db.Where("spot_id in ?", spotIDs).Find(&reservations)
	if reservationResult.Error != nil || len(reservations) == 0 {
		c.String(http.StatusNotFound, "Couldn't find the reservations associated with the organization")
	}

	reservationMap := make(map[uint][]database.Reservation)
	for _, reservation := range reservations {
		reservationMap[reservation.SpotID] = append(reservationMap[reservation.SpotID], reservation)
	}

	for i := range spots {
		if res, found := reservationMap[spots[i].ID]; found {
			spots[i].Reservations = res
		}
	}

	/*
		// Queries to gather the spots attached to the organization
		// and their reservations
		spots := database.Spot{}
		spotResult := database.Db.Model(&database.Spot{}).Preload("Reservations").Where("spotID IN ?", spotIDs).Find(&spots)
		if spotResult.Error != nil {
			c.String(http.StatusNotFound, "Couldn't find the spots associated with you")
			return
		}
	*/

	c.IndentedJSON(http.StatusOK, spots) // Send back 220 with the JSON of the spots & reservations

}
