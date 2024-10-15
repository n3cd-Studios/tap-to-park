package routes

import (
	"net/http"
	"strconv"
	"tap-to-park/database"
	"time"

	"github.com/gin-gonic/gin"
)

type OrganizationRoutes struct{}

// GetOrganization godoc
// @Summary      Get all of the organizations associated with an admin
// @Produce      json
// @Success      200  {object}  database.Organization
// @Failure      400  {string}  "Unauthorized"
// @Router       /organization/me [get]
func (*OrganizationRoutes) GetOrganization(c *gin.Context) {

	uuid := c.MustGet("guid")

	user := database.User{}
	if result := database.Db.Where("guid = ?", uuid).First(&user); result.Error != nil {
		c.String(http.StatusNotFound, "For some reason, you don't exist!")
		return
	}

	deep := c.Query("deep")

	result := database.Db.Model(&database.Organization{}).Where("id = ?", user.OrganizationID)
	if deep == "true" {
		result = result.Preload("Spots.Reservations")
	}

	organization := database.Organization{}
	result = result.First(&organization)
	if result.Error != nil {
		c.String(http.StatusNotFound, "Couldn't find the organization associated with you")
		return
	}

	c.IndentedJSON(http.StatusOK, organization)
}

type GetSpotsOutput struct {
	Items []database.Spot `json:"items"`
	Pages int64           `json:"pages"`
	Page  int64           `json:"page"`
}

// GetSpots godoc
// @Summary      Get all of the spots associated with an organization
// @Produce      json
// @Success      200  {array}   []database.Spot
// @Failure      400  {string}  "Unauthorized"
// @Router       /organization/spots [get]
func (*OrganizationRoutes) GetSpots(c *gin.Context) {

	uuid := c.MustGet("guid")

	user := database.User{}
	if result := database.Db.Where("guid = ?", uuid).First(&user); result.Error != nil {
		c.String(http.StatusNotFound, "For some reason, you don't exist!")
		return
	}

	page, perr := strconv.ParseInt(c.Query("page"), 10, 64)
	if perr != nil {
		page = 0
	}

	size, perr := strconv.ParseInt(c.Query("size"), 10, 64)
	if perr != nil {
		size = 10
	}

	spots := []database.Spot{}
	count := int64(0)
	result := database.Db.Model(&database.Spot{}).Where("organization_id = ?", user.OrganizationID).Count(&count).Offset(int(size * page)).Limit(int(size)).Find(&spots)
	if result.Error != nil {
		c.String(http.StatusConflict, "Couldn't count all of the spots in the organization.")
		return
	}

	c.IndentedJSON(http.StatusOK, GetSpotsOutput{
		Items: spots,
		Pages: (count / size),
		Page:  page,
	})
}

// CreateInvite godoc
// @Summary      Create an invite to allow new user to join admin's organization
// @Produce      json
// @Success      200  {object}  database.Invite
// @Failure      401  {string}  "Unauthorized"
// @Failure      404  {string}  "User or Organization not found"
// @Failure      500  {string}  "Failed to create invite"
// @Router       /organization/code [post]
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

// GetInvites godoc
// @Summary      Get's all the invites for an organization
// @Produce      json
// @Success      200  {object}  database.Invite
// @Failure      401  {string}  "Unauthorized"
// @Failure      404  {string}  "User or Organization not found"
// @Failure      500  {string}  "Failed to create invite"
// @Router       /organization/code [post]
// @Security     BearerAuth
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
