package main

import (
	"net/http"
	"tap-to-park/database"

	"github.com/gin-gonic/gin"
)

type Service struct{}

func (*Service) getReservationByID(c *gin.Context) {

	id := c.Param("id")

	var reservation database.Reservation
	result := database.Db.Where("id = ?", id).First(&reservation)
	err := result.Error

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusAccepted, reservation)
}

type ReservationInput struct {
	SpotID uint `json:"spotID"`
}

func (*Service) postReservation(c *gin.Context) {
	var input ReservationInput
	if err := c.BindJSON(&input); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	reservation := database.Reservation{SpotID: input.SpotID}
	result := database.Db.Create(&reservation)
	if err := result.Error; err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, reservation)
}
