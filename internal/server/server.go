package server

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/acernik/go-server-boilerplate/internal/service"
)

// Server holds the Server data.
type Server struct {
	Service    service.Service
	Logger     *zap.Logger
	TimeFormat string
}

// New returns new server.
func New(s service.Service, logger *zap.Logger, timeFormat string) *Server {
	return &Server{
		Service:    s,
		Logger:     logger,
		TimeFormat: timeFormat,
	}
}

// RegisterRoutes registers all the  API routes.
func (s *Server) RegisterRoutes() *gin.Engine {
	router := gin.Default()

	sr := router.Group("/api/v1/test-items")

	sr.POST("/create", s.Insert)
	sr.GET("/health", s.Health)

	return router
}
