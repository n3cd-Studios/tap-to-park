package routes

import (
	"net/http"
	"strconv"
	"tap-to-park/database"
	"time"

	"github.com/gin-gonic/gin"
)

type ReservationRoutes struct{}

// GetReservation godoc
//
// @Summary		Get a reservation by ID
// @Description	Get a reservation for a Spot based on the Reservation's GUID
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

type CreateFakeReservationInput struct {
	SpotId 	 string 	`json:"spot_id" binding:"required"`
	Email    string 	`json:"email" binding:"required"`
	Start 	 int64 	`json:"start" binding:"required"`
	Minutes  int 		`json:"minutes" binding:"required"`
	Cost     float64	`json:"cost" binding:"required"`
}

// This function cannot be used in production anyways!
func (*ReservationRoutes) CreateFakeReservation(c *gin.Context) {

	input := CreateFakeReservationInput{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.String(http.StatusBadRequest, "Invalid input.")
		return
	}

	spot := database.Spot{}
	if result := database.Db.Where("guid = ?", input.SpotId).First(&spot); result.Error != nil {
		c.String(http.StatusBadRequest, "That spot ID is invalid.")
		return
	}

	now := time.UnixMilli(input.Start);
	reservation := database.Reservation{
		Start:               now,
		End:                 now.Add(time.Duration(input.Minutes*60) * time.Minute),
		Email:               input.Email,
		Price:               input.Cost,
		SpotID:              spot.ID,
		StripeTransactionID: "FAKE-" + input.Email + "-" + strconv.FormatInt(now.UnixMilli(), 10),
	}
	if result := database.Db.Create(&reservation); result.Error != nil {
		c.String(http.StatusBadRequest, "Something went wrong!")
		return
	}

	c.String(http.StatusOK, reservation.Guid);
}
