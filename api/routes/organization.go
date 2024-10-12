package routes

import (
	"net/http"
	"tap-to-park/database"
	"time"

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
	result := database.Db.Preload("Spots.Reservations").Where("id = ?", user.OrganizationID).First(&organization)
	if result.Error != nil {
		c.String(http.StatusNotFound, "Couldn't find the organization associated with you")
		return
	}

	c.IndentedJSON(http.StatusOK, &organization) // Send back 220 with the JSON of the spots & reservations

}

// CreateInvite godoc
// @Summary      Create an invite to allow new user to join admin's organization
// @Produce      json
// @Success      200  {object}  database.Invite
// @Failure      401  {string}  "Unauthorized"
// @Failure      404  {string}  "User or Organization not found"
// @Failure      500  {string}  "Failed to create invite"
// @Router       /admin/organization [post]
// @Security     BearerAuth
func (*OrganizationRoutes) CreateInvite(c *gin.Context) {

	uuid := c.MustGet("guid")

	user := database.User{}
	if result := database.Db.Where("guid = ?", uuid).First(&user); result.Error != nil {
		c.String(http.StatusNotFound, "For some reason, you don't exist!")
		return
	}

	organization := database.Organization{}
	result := database.Db.Model(&database.Organization{}).Where("id = ?", user.OrganizationID).First(&organization)
	if result.Error != nil {
		c.String(http.StatusNotFound, "Couldn't find the organization associated with you")
		return
	}

	invite := database.Invite{Expiration: time.Now().Add(1 * time.Hour), OrganizationID: organization.ID, CreatedByID: user.ID}

	maxGenerationAttempts := 3
	for attempts := 0; attempts < maxGenerationAttempts; attempts++ {
		err := database.Db.Create(&invite).Error
		if err == nil {
			c.IndentedJSON(http.StatusOK, invite)
			return
		}
	}

	// After failed attempts, return an error
	c.String(http.StatusInternalServerError, "Failed to create invite.")
}

func (*OrganizationRoutes) GetInvites(c *gin.Context) {
	uuid := c.MustGet("guid")

	user := database.User{}
	if result := database.Db.Where("guid = ?", uuid).First(&user); result.Error != nil {
		c.String(http.StatusNotFound, "You don't exist lol")
		return
	}

	organization := database.Organization{}
	if result := database.Db.Where("id = ?", user.OrganizationID).First(&organization); result.Error != nil {
		c.String(http.StatusNotFound, "Could not find your organization")
		return
	}

	var invites []database.Invite
	if result := database.Db.Where("organization_id = ?", organization.ID).Find(&invites); result.Error != nil {
		c.String(http.StatusNotFound, "Could not find any invites for your organization")
		return
	}

	c.IndentedJSON(http.StatusOK, invites)
}
