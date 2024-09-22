package routes

import (
	"net/http"
	"strconv"
	"tap-to-park/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

type SpotRoutes struct{}

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

	point := database.Point{X: lat, Y: lng}

	var spots []database.Spot
	result := database.Db.Order(clause.OrderBy{Expression: clause.Expr{SQL: "coords <-> Point ?", Vars: []interface{}{[]database.Point{point}}, WithoutParentheses: true}}).Limit(10).Find(&spots)
	err := result.Error

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusAccepted, spots)
}
