package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/acernik/go-server-boilerplate/internal/server/mapper"
	"github.com/acernik/go-server-boilerplate/internal/server/model"
)

func isRequestValid(req model.InsertTestItemRequest) bool {
	if len(req.Name) == 0 || len(req.Description) == 0 {
		return false
	}

	return true
}

func (s *Server) Insert(c *gin.Context) {
	var req model.InsertTestItemRequest

	err := c.BindJSON(&req)
	if err != nil {
		s.Logger.Error(
			"failed to bind JSON to value of type InsertTestItemRequest",
			zap.String("err", err.Error()),
			zap.String("time", time.Now().UTC().Format(s.TimeFormat)),
		)

		c.JSON(http.StatusBadRequest, gin.H{})

		return
	}

	if !isRequestValid(req) {
		c.JSON(http.StatusBadRequest, gin.H{})
	}

	_, err = s.Service.InsertTestTableItem(mapper.MapTestTableItemAPIToService(req))
	if err != nil {
		s.Logger.Error(
			"failed while executing Service.InsertTestTableItem",
			zap.String("err", err.Error()),
			zap.String("time", time.Now().UTC().Format(s.TimeFormat)),
		)
	}

	c.JSON(http.StatusCreated, gin.H{})
}
