package auth

import (
	"net/http"
	"tap-to-park/auth"
	"tap-to-park/database"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session_id, err := auth.Get(c.Request.Header.Get("Authentication"))
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized.")
			c.Abort()
			return
		}

		session := database.Session{}
		if result := database.Db.Where("guid = ?", session_id).Where("expires > ?", time.Now()).First(&session); result.Error != nil {
			c.String(http.StatusUnauthorized, "Unauthorized.")
			c.Abort()
			return
		}

		session.LastUsed = time.Now()
		database.Db.Updates(&session)

		user := database.User{}
		if result := database.Db.Where("id = ?", session.UserID).First(&user); result.Error != nil {
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
// @Router		/auth/login [post]
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
	request := c.Request
	token, err := auth.Generate(user.ID, request.UserAgent(), request.Host)
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
// @Router		/auth/register [post]
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

	request := c.Request
	token, err := auth.Generate(user.ID, request.UserAgent(), request.Host)
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

type OAuthLogin interface {
	Initialize(c *gin.Context)
	Callback(c *gin.Context) *database.User
}

var oauthTypes = map[string]OAuthLogin{"github": Github{}}

func (*AuthRoutes) OAuthInitialize(c *gin.Context) {

	authType := c.Param("type")
	handler, exists := oauthTypes[authType]
	if !exists {
		c.String(http.StatusBadRequest, "That OAuth flow does not exist.")
		return
	}

	handler.Initialize(c)
}

func (*AuthRoutes) OAuthCallback(c *gin.Context) {

	authType := c.Param("type")
	handler, exists := oauthTypes[authType]
	if !exists {
		c.String(http.StatusBadRequest, "That OAuth flow does not exist.")
		return
	}

	user := handler.Callback(c)
	if user == nil {
		c.String(http.StatusBadRequest, "OAuth flow failed to sign you in.")
		return
	}

	request := c.Request
	token, err := auth.Generate(user.ID, request.UserAgent(), request.Host)
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to create session.")
		return
	}

	c.IndentedJSON(http.StatusOK, JWTResponse{Token: token})
}

// GetInfo godoc
//
// @Summary		Get user info
// @Description Get a user's info based on a Bearer token
// @Tags		auth
// @Accept		json
// @Produce		json
// @Success		200	{object}	database.User
// @Failure		401	{string} string "Unauthorized."
// @Router		/auth/info [post]
// @Security	BearerToken
func (*AuthRoutes) GetInfo(c *gin.Context) {
	user := c.MustGet("user").(database.User)
	c.IndentedJSON(http.StatusOK, user)
}

// GetSessions godoc
//
// @Summary		Get user's sessions
// @Description Get a user's sessions based on a Bearer token
// @Tags		auth
// @Accept		json
// @Produce		json
// @Success		200	{array}	[]database.Session
// @Failure		401	{string} string "Unauthorized."
// @Failure		404	{string} string "You don't have any sessions."
// @Router		/auth/sessions [get]
// @Security	BearerToken
func (*AuthRoutes) GetSessions(c *gin.Context) {
	user := c.MustGet("user").(database.User)

	sessions := []database.Session{}
	if result := database.Db.Where("user_id = ?", user.ID).Find(&sessions); result.Error != nil {
		c.String(http.StatusNotFound, "You don't have any sessions.")
		return
	}

	c.IndentedJSON(http.StatusOK, sessions)
}

// RevokeSession godoc
//
// @Summary		Revoke a session
// @Description Revoke a session based on an ID
// @Tags		auth
// @Accept		json
// @Produce		json
// @Success		200	{array}	string "Revoked session."
// @Failure		401	{string} string "Unauthorized."
// @Failure		404	{string} string "Failed to revoke session."
// @Router		/auth/sessions/{id} [delete]
// @Security	BearerToken
func (*AuthRoutes) RevokeSession(c *gin.Context) {
	user := c.MustGet("user").(database.User)
	id := c.Param("id")

	session := &database.Session{}
	if result := database.Db.Where("guid = ?", id).Where("user_id = ?", user.ID).First(&session).Delete(&session); result.Error != nil {
		c.String(http.StatusNotFound, "Failed to revoke session.")
		return
	}

	c.String(http.StatusOK, "Revoked session.")
}
