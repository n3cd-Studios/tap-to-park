package routes

import (
	"net/http"
	"strconv"
	"tap-to-park/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

type SpotRoutes struct{}

// GetSpotsNear godoc
// @Summary      Get the spots near a longitude and latitude
// @Produce      json
// @Param        lat    query     number  true  "latitude to search by"
// @Param        lng    query     number  true  "longitude to search by"
// @Success      200  {object}  database.Spot
// @Failure      404  {object}  database.Error
// @Router       /spots/near [get]
func (*SpotRoutes) GetSpotsNear(c *gin.Context) {

	latParam := c.Query("lat")
	lngParam := c.Query("lng")

	lat, perr := strconv.ParseFloat(latParam, 64)
	if perr != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Latitude must me a float."})
		return
	}

	lng, perr := strconv.ParseFloat(lngParam, 64)
	if perr != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Longitude must me a float."})
		return
	}

	point := database.Coordinates{Longitude: lng, Latitude: lat}

	spots := []database.Spot{}
	result := database.Db.Order(clause.OrderBy{Expression: clause.Expr{SQL: "coords <-> Point (?,?)", Vars: []interface{}{point.Latitude, point.Longitude}, WithoutParentheses: true}}).Limit(10).Find(&spots)
	err := result.Error

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusAccepted, spots)
}
