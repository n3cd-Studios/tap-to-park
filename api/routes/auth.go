package routes

import (
	"net/http"
	"tap-to-park/auth"
	"tap-to-park/database"

	"github.com/gin-gonic/gin"
)

type AuthRoutes struct{}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := auth.TokenValid(c)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login godoc
// @Summary      Logs a User in using a username and a password
// @Produce      json
// @Success      200  {object}  LoginInput
// @Failure      404  {object}  database.Error
// @Router       /auth/login [post]
func (*AuthRoutes) Login(c *gin.Context) {

	var input LoginInput
	if err := c.BindJSON(&input); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid body recieved."})
		return
	}

	user := database.User{}
	result := database.Db.Where("email = ?", input.Email).First(&user)
	if result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Failed to login."})
		return
	}

	match, err := auth.ComparePasswordAndHash(input.Password, user.PasswordHash)
	if !match || err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Failed to login."})
		return
	}

	token, err := auth.GenerateToken(user.UniqueID)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Failed to login."})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"token": token})
}

type RegisterInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register godoc
// @Summary      Registers a User in using a username and a password
// @Produce      json
// @Success      200  {object}  RegisterInput
// @Failure      404  {object}  database.Error
// @Router       /auth/register [post]
func (*AuthRoutes) Register(c *gin.Context) {

	// TODO: CHANGE ALL ERRORS TO GENERIC ERROR FOR SECURITY

	var input RegisterInput
	if err := c.BindJSON(&input); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid body recieved."})
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
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Hash failed"})
		return
	}

	user := database.User{Email: input.Email, PasswordHash: hash, OrganizationID: 1}
	if err := database.Db.Create(&user); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Creation failed"})
		return
	}

	token, err := auth.GenerateToken(user.UniqueID)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Token failed"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"token": token})
}
