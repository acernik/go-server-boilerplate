package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) Health(c *gin.Context) {
	c.Status(http.StatusOK)
}
