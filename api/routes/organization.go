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
