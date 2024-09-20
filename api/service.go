package api

import (
	"tap-to-park/api/database"

	"github.com/gin-gonic/gin"
)

type Service struct{}

func (*Service) getReservationByID(c *gin.Context) {

	id := c.Param("id")

	var reservation database.Reservation
	result := database.Db.Where("id = ?", id).First(&reservation)
	err := result.Error

	if err != nil {
		c.IndentedJSON(400, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(200, reservation)
}
