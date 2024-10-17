package routes

import (
	"net/http"
	"tap-to-park/database"

	"github.com/gin-gonic/gin"
)

type ReservationRoutes struct{}

// GetReservation godoc
//
// @Summary		Create an invite
// @Description	Create an invite for User's organization based on their Bearer token
// @Tags		reservation
// @Accept		json
// @Produce		json
// @Success		200	{object} database.Reservation
// @Failure		404 {string} string "That reservation does not exist."
// @Router		/reservation/{id} [get]
func (*ReservationRoutes) GetReservation(c *gin.Context) {

	id := c.Param("id")
	reservation := database.Reservation{}
	if result := database.Db.Where("id = ?", id).First(&reservation); result.Error != nil {
		c.IndentedJSON(http.StatusNotFound, "That reservation does not exist.")
		return
	}

	c.IndentedJSON(http.StatusOK, reservation)
}
