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

	// Queries to gather the spots attached to the organization
	// and their reservations
	spots := database.Spot{}
	spotResult := database.Db.Model(&database.Spot{}).Preload("Reservations").Where("spotID IN ?", spotIDs).Find(&spots)
	if spotResult.Error != nil {
		c.String(http.StatusNotFound, "Couldn't find the spots associated with you")
		return
	}

	c.IndentedJSON(http.StatusOK, spots) // Send back 220 with the JSON of the spots & reservations

}
