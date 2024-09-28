package routes

import (
	"math/rand"
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

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// thx https://www.calhoun.io/creating-random-strings-in-go/
var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func GenerateAlphanumeric(length int, charset string) string {
	code := make([]byte, length)
	for i := range code {
		code[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(code)
}

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

	const maxGenerationAttempts = 10
	var inviteCode string
	for attempts := 0; attempts < maxGenerationAttempts; attempts++ {
		inviteCode = GenerateAlphanumeric(9, charset)
		var existingInvite database.Invite
		if err := database.Db.Where("code = ?", inviteCode).First(&existingInvite).Error; err != nil {
			// break if code has not been used
			break
		}
	}
	invite := database.Invite{ID: GenerateAlphanumeric(9, charset), Start: time.Now(), End: time.Now().Add(1 * time.Hour), OrganizationID: organization.ID, CreatedByID: user.ID}

	if err := database.Db.Create(&invite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create invite."})
		return
	}

	c.IndentedJSON(http.StatusOK, invite)
}
