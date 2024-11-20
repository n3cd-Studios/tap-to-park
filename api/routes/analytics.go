package routes

import (
	"net/http"
	"tap-to-park/database"

	"github.com/gin-gonic/gin"
)

type AnalyticRoutes struct{}

// GetTopSpots godoc
//
// @Summary		Get top spots
// @Description	Get the top spots for an organization associated with a User based on their Bearer token
// @Tags		organization,analytics
// @Accept		json
// @Produce		json
// @Success		200	{array} []map[string]interface
// @Failure		400	{string} string	"You don't seem to have an organization."
// @Failure		401	{string} string "Unauthorized."
// @Failure		404	{string} string "Failed to generate analytic."
// @Router		/analytics/top [get]
// @Security 	BearerToken
func (*AnalyticRoutes) GetTopSpots(c *gin.Context) {

	user := c.MustGet("user").(database.User)

	organization := database.Organization{}
	result := database.Db.Model(&database.Organization{}).Where("id = ?", user.OrganizationID).First(&organization)
	if result.Error != nil {
		c.String(http.StatusNotFound, "You don't seem to have an organization.")
		return
	}

	// Silly you have to do this to get this to work
	spots := []map[string]interface{}{}
	result = database.Db.Model(&database.Spot{}).
		Select("spots.name name", "spots.id id", "ROUND(SUM(price)/100.0) revenue").
		Joins("JOIN reservations ON spots.id = reservations.spot_id").
		Group("spots.id, spots.name").
		Order("revenue DESC").
		Where("organization_id = ?", organization.ID).
		Limit(10).
		Scan(&spots)
	if result.Error != nil {
		c.String(http.StatusNotFound, "Failed to generate analytic.")
		return
	}

	c.IndentedJSON(http.StatusOK, spots)
}

// GetPeakTimes godoc
//
// @Summary		Get peak times
// @Description	Get the peak times for spots in an organization associated with a User based on their Bearer token
// @Tags		organization,analytics
// @Accept		json
// @Produce		json
// @Success		200	{array} []map[string]interface
// @Failure		400	{string} string	"You don't seem to have an organization."
// @Failure		401	{string} string "Unauthorized."
// @Failure		404	{string} string "Failed to generate analytic."
// @Router		/analytics/peak [get]
// @Security 	BearerToken
func (*AnalyticRoutes) GetPeakTimes(c *gin.Context) {

	user := c.MustGet("user").(database.User)

	organization := database.Organization{}
	result := database.Db.Model(&database.Organization{}).Where("id = ?", user.OrganizationID).First(&organization)
	if result.Error != nil {
		c.String(http.StatusNotFound, "You don't seem to have an organization.")
		return
	}

	// Silly you have to do this to get this to work
	times := []map[string]interface{}{}
	result = database.Db.Model(&database.Spot{}).
		Select("COUNT(*) amount", "EXTRACT(hour from start::timestamp) time", "SUM(price) revenue").
		Joins("JOIN reservations ON spots.id = reservations.spot_id").
		Group("time").
		Order("time ASC").
		Where("organization_id = ?", organization.ID).
		Scan(&times)
	if result.Error != nil {
		c.String(http.StatusNotFound, "Failed to generate analytic.")
		return
	}

	c.IndentedJSON(http.StatusOK, times)
}

// GetRevenueByMonth godoc
//
// @Summary		Get revenue by month
// @Description	Get the revenue by month from an organization associated with a User based on their Bearer token
// @Tags		organization,analytics
// @Accept		json
// @Produce		json
// @Success		200	{array} []map[string]interface
// @Failure		400	{string} string	"You don't seem to have an organization."
// @Failure		401	{string} string "Unauthorized."
// @Failure		404	{string} string "Failed to generate analytic."
// @Router		/analytics/revenue [get]
// @Security 	BearerToken
func (*AnalyticRoutes) GetRevenueByMonth(c *gin.Context) {

	user := c.MustGet("user").(database.User)

	organization := database.Organization{}
	result := database.Db.Model(&database.Organization{}).Where("id = ?", user.OrganizationID).First(&organization)
	if result.Error != nil {
		c.String(http.StatusNotFound, "You don't seem to have an organization.")
		return
	}

	// Silly you have to do this to get this to work
	months := []map[string]interface{}{}
	result = database.Db.Model(&database.Spot{}).
		Select("SUM(price) revenue", "(EXTRACT(month from start::timestamp)) AS month").
		Joins("JOIN reservations ON spots.id = reservations.spot_id").
		Group("month").
		Order("month ASC").
		Where("organization_id = ?", organization.ID).
		Scan(&months)
	if result.Error != nil {
		c.String(http.StatusNotFound, "Failed to generate analytic.")
		return
	}

	c.IndentedJSON(http.StatusOK, months)
}
