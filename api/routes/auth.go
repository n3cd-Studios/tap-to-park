package routes

import (
	"net/http"
	"tap-to-park/auth"
	"tap-to-park/database"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := auth.TokenExtract(c.Request.Header.Get("Authentication"))
		guid, err := auth.TokenExtractID(token)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		c.Set("guid", guid)
		c.Next()
	}
}

type AuthRoutes struct{}

type JWTResponse struct {
	Token string `json:"token"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login godoc
// @Summary      Logs a User in using a username and a password
// @Produce      json
// @Success      200  {object}  JWTResponse
// @Failure      400  {string}  "Failed to log in"
// @Router       /auth/login [post]
func (*AuthRoutes) Login(c *gin.Context) {

	var input LoginInput
	if err := c.BindJSON(&input); err != nil {
		c.String(http.StatusBadRequest, "Invalid body recieved.")
		return
	}

	user := database.User{}
	result := database.Db.Where("email = ?", input.Email).First(&user)
	if result.Error != nil {
		c.String(http.StatusBadRequest, "Failed to login.")
		return
	}

	match, err := auth.ComparePasswordAndHash(input.Password, user.PasswordHash)
	if !match || err != nil {
		c.String(http.StatusBadRequest, "Failed to login.")
		return
	}

	token, err := auth.TokenGenerate(user.Guid)
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to login.")
		return
	}

	c.IndentedJSON(http.StatusOK, JWTResponse{Token: token})
}

type RegisterInput struct {
	Email      string `json:"email" binding:"required"`
	Password   string `json:"password" binding:"required"`
	InviteCode string `json:"invite"`
}

// Register godoc
// @Summary      Registers a User using a username, password, and an optional invite code
// @Produce      json
// @Success      200  {object}  JWTResponse
// @Failure      400  {string}  "Failed to register user"
// @Router       /auth/register [post]
func (*AuthRoutes) Register(c *gin.Context) {

	// TODO: CHANGE ALL ERRORS TO GENERIC ERROR FOR SECURITY

	var input RegisterInput
	if err := c.BindJSON(&input); err != nil {
		c.String(http.StatusBadRequest, "Invalid body recieved.")
		return
	}

	var existingUser database.User
	if result := database.Db.Where("email = ?", input.Email).First(&existingUser); result.Error == nil {
		c.String(http.StatusBadRequest, "Email already in use.")
		return
	}

	var organizationID uint
	if input.InviteCode != "" {
		// Validate invite code
		var invite database.Invite
		if result := database.Db.Where("ID = ?", input.InviteCode).First(&invite); result.Error != nil {
			c.String(http.StatusBadRequest, "Invalid invite code.")
			return
		}

		if time.Now().After(invite.Expiration) || invite.UsedByID != 0 {
			c.String(http.StatusBadRequest, "Invalid or expired invite code.")
			return
		}
		organizationID = invite.OrganizationID
	} else {
		// No invite code
		c.String(http.StatusBadRequest, "Currently, you need an invite to register.")
		organizationID = 0
	}

	hash, err := auth.GenerateFromPassword(input.Password, &auth.Params{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	})

	if err != nil {
		c.String(http.StatusBadRequest, "Failed to create user.")
		return
	}

	user := database.User{Email: input.Email, PasswordHash: hash, OrganizationID: organizationID}

	if err := database.Db.Create(&user).Error; err != nil {
		c.String(http.StatusBadRequest, "Failed to create user.")
		return
	}

	// mark inv code used
	if input.InviteCode != "" {
		var invite database.Invite
		invite.UsedByID = user.ID
		if err := database.Db.Save(&invite).Error; err != nil {
			c.String(http.StatusInternalServerError, "Failed to update invite.")
			return
		}
	}

	token, err := auth.TokenGenerate(user.Guid)
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to create user.")
		return
	}

	c.IndentedJSON(http.StatusOK, JWTResponse{Token: token})
}

// Info godoc
// @Summary      Gets the info of the current user
// @Produce      json
// @Success      200  {object}  JWTResponse
// @Failure      400  {string}  "Failed to use token to retrieve user information"
// @Router       /auth/info [get]
func (*AuthRoutes) Info(c *gin.Context) {

	uuid := c.MustGet("guid")

	user := database.User{}
	if result := database.Db.Where("guid = ?", uuid).First(&user); result.Error != nil {
		c.String(http.StatusNotFound, "For some reason, you don't exist!")
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}
