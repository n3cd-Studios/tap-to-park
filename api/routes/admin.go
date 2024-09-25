package routes

import (
	"net/http"
	"tap-to-park/database"

	"github.com/gin-gonic/gin"
)

type AdminRoutes struct{}

// GetOrganizations godoc
// @Summary      Get all of the organizations associated with an admin
// @Produce      json
// @Success      200  {array} []database.Organization
// @Failure      400  {string}  "Unauthorized"
// @Router       /admin/organization [get]
func (*AdminRoutes) GetOrganization(c *gin.Context) {

	uuid := c.MustGet("uuid")

	user := database.User{}
	if result := database.Db.Where("unique_id = ?", uuid).First(&user); result.Error != nil {
		c.String(http.StatusNotFound, "For some reason, you don't exist!")
		return
	}

	organization := database.Organization{}
	if result := database.Db.Where("id = ?", user.OrganizationID).First(&organization); result.Error != nil {
		c.String(http.StatusNotFound, "Couldn't find organizations associated with you.")
		return
	}

	c.IndentedJSON(http.StatusOK, organization)
}
