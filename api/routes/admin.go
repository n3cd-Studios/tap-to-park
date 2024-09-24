package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminRoutes struct{}

// Login godoc
// @Summary      Logs a User in using a username and a password
// @Produce      json
// @Success      200  {string}  "Logged in"
// @Failure      400  {string}  "Unauthorized"
// @Router       /auth/test [get]
func (*AdminRoutes) Test(c *gin.Context) {
	c.String(http.StatusOK, "Logged in")
}
