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
//
// @Summary		Get your organization
// @Description	Get the organization associated with a User based on their Bearer token
// @Tags		organization
// @Accept		json
// @Produce		json
// @Param		deep  query		boolean	false	"Pull a deep copy of all of the organization's information"
// @Success		200	{object} database.Organization
// @Failure		400	{string} string	"You don't seem to have an organization."
// @Failure		401	{string} string "Unauthorized."
// @Router		/organization/me [get]
// @Security 	BearerToken
func (*OrganizationRoutes) GetOrganization(c *gin.Context) {

	user := c.MustGet("user").(database.User)
	deep := c.Query("deep")

	result := database.Db.Model(&database.Organization{}).Where("id = ?", user.OrganizationID)
	if deep == "true" {
		result = result.Preload("Spots.Reservations")
	}

	organization := database.Organization{}
	result = result.First(&organization)
	if result.Error != nil {
		c.String(http.StatusNotFound, "You don't seem to have an organization.")
		return
	}

	c.IndentedJSON(http.StatusOK, organization)
}

// GetSpots godoc
//
// @Summary		Get the spots for your organization
// @Description	Get the spots associated with a User's organization based on their Bearer token
// @Tags		organization,spot
// @Accept		json
// @Produce		json
// @Param		size  query		number	false	"The size of a page"
// @Param		page  query		number	false	"The page"
// @Success		200	{array} []database.Spot
// @Failure		500	{string} string	"Couldn't count all of the spots in the organization."
// @Failure		401	{string} string "Unauthorized."
// @Router		/organization/spots [get]
// @Security 	BearerToken
func (*OrganizationRoutes) GetSpots(c *gin.Context) {

	user := c.MustGet("user").(database.User)

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
		c.String(http.StatusInternalServerError, "Couldn't count all of the spots in the organization.")
		return
	}

	c.IndentedJSON(http.StatusOK, PaginatorOutput[database.Spot]{
		Items: spots,
		Pages: (count / size),
		Page:  page,
	})
}

// CreateInvite godoc
//
// @Summary		Create an invite
// @Description	Create an invite for User's organization based on their Bearer token
// @Tags		organization,invite
// @Accept		json
// @Produce		json
// @Success		200	{object} database.Invite
// @Failure		404 {string} string "Failed to find your organization."
// @Failure		500	{string} string	"Failed to create invite."
// @Failure		401	{string} string "Unauthorized."
// @Router		/organization/invites [post]
// @Security 	BearerToken
func (*OrganizationRoutes) CreateInvite(c *gin.Context) {

	user := c.MustGet("user").(database.User)
	organization := database.Organization{}
	result := database.Db.Model(&database.Organization{}).Where("id = ?", user.OrganizationID).First(&organization)
	if result.Error != nil {
		c.String(http.StatusNotFound, "Failed to find your organization.")
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
//
// @Summary		Get the invites for your organization
// @Description	Get the invites associated with a User's organization based on their Bearer token
// @Tags		organization,invite
// @Accept		json
// @Produce		json
// @Param		size  query		number	false	"The size of a page"
// @Param		page  query		number	false	"The page"
// @Success		200	{array} []database.Invite
// @Failure		404 {string} string	"No invites were found for your organization."
// @Failure		500	{string} string	"Couldn't count all of the invites in the organization."
// @Failure		401	{string} string "Unauthorized."
// @Router		/organization/invites [get]
// @Security 	BearerToken
func (*OrganizationRoutes) GetInvites(c *gin.Context) {
	user := c.MustGet("user").(database.User)

	page, perr := strconv.ParseInt(c.Query("page"), 10, 64)
	if perr != nil {
		page = 0
	}

	size, perr := strconv.ParseInt(c.Query("size"), 10, 64)
	if perr != nil {
		size = 10
	}

	organization := database.Organization{}
	if result := database.Db.Where("id = ?", user.OrganizationID).First(&organization); result.Error != nil {
		c.String(http.StatusNotFound, "You don't seem to have an organization.")
		return
	}

	count := int64(0)
	var invites []database.Invite
	if result := database.Db.Where("organization_id = ?", organization.ID).Count(&count).Offset(int(page * size)).Limit(int(size)).Find(&invites); result.Error != nil {
		c.String(http.StatusNotFound, "No invites were found for your organization.")
		return
	}

	c.IndentedJSON(http.StatusOK, PaginatorOutput[database.Invite]{
		Items: invites,
		Pages: (count / size),
		Page:  page,
	})
}
