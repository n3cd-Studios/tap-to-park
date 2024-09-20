package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Service struct {
}

// getAlbums responds with the list of all albums as JSON.
func (*Service) getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{})
}
