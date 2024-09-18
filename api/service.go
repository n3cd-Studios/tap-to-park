package api

import (
	"github.com/gin-gonic/gin"
)

type Service struct {
}

// getAlbums responds with the list of all albums as JSON.
func (*Service) getAlbums(c *gin.Context) {
	c.Done()
}
