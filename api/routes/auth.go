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
			c.String(http.StatusUnauthorized, "Unauthorized.")
			c.Abort()
			return
		}
		user := database.User{}
		if result := database.Db.Where("guid = ?", guid).Find(&user); result.Error != nil {
			c.String(http.StatusUnauthorized, "Unauthorized.")
			c.Abort()
			return
		}
		c.Set("user", user)
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
//
// @Summary		Login a user
// @Description Login a user, this will generate a Bearer token to be used with Authenticated requests.
// @Tags		auth
// @Accept		json
// @Produce		json
// @Success		200	{object}	JWTResponse
// @Failure		400	{string}	string	"Failed to login."
// @Failure		400	{string}	string	"Invalid body recieved."
// @Router		/login [post]
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
//
// @Summary		Register a user
// @Description Register a user using an organization's invite code, this will generate a Bearer token to be used with Authenticated requests.
// @Tags		auth
// @Accept		json
// @Produce		json
// @Success		200	{object}	JWTResponse
// @Failure		400	{string}	string	"Failed to register."
// @Failure		400	{string}	string	"Invalid body recieved."
// @Failure		500	{string}	string	"Failed to update invite."
// @Router		/register [post]
func (*AuthRoutes) Register(c *gin.Context) {

	// TODO: CHANGE ALL ERRORS TO GENERIC ERROR FOR SECURITY

	var input RegisterInput
	if err := c.BindJSON(&input); err != nil {
		c.String(http.StatusBadRequest, "Invalid body recieved.")
		return
	}

	hash, err := auth.GenerateFromPassword(input.Password, &auth.Params{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	})

	if err != nil {
		c.String(http.StatusBadRequest, "Failed to register.")
		return
	}

	var invite database.Invite
	var organizationID uint
	if result := database.Db.Where("ID = ?", input.InviteCode).First(&invite); result.Error != nil {
		c.String(http.StatusBadRequest, "Failed to register.")
		return
	}

	if time.Now().After(invite.Expiration) || invite.UsedByID != 0 {
		c.String(http.StatusBadRequest, "Failed to register.")
		return
	}

	organizationID = invite.OrganizationID

	user := database.User{Email: input.Email, PasswordHash: hash, OrganizationID: organizationID}

	if err := database.Db.Create(&user).Error; err != nil {
		c.String(http.StatusBadRequest, "Failed to register.")
		return
	}

	token, err := auth.TokenGenerate(user.Guid)
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to create user.")
		return
	}

	invite.UsedByID = user.ID
	if err := database.Db.Save(&invite).Error; err != nil {
		c.String(http.StatusInternalServerError, "Failed to update invite.")
		return
	}

	c.IndentedJSON(http.StatusOK, JWTResponse{Token: token})
}

// Info godoc
//
// @Summary		Get user info
// @Description Get a user's info based on a Bearer token
// @Tags		auth
// @Accept		json
// @Produce		json
// @Success		200	{object}	database.User
// @Failure		401	{string} string "Unauthorized."
// @Router		/info [post]
// @Security	BearerToken
func (*AuthRoutes) Info(c *gin.Context) {
	user := c.MustGet("user").(database.User)
	c.IndentedJSON(http.StatusOK, user)
}
