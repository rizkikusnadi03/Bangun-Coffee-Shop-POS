package handler

import (
	"strings"

	"github.com/gin-gonic/gin"
	"web-coffee/pkg/response"
)

// NewRouter menginisialisasi router Gin dengan middleware, mode environment, dan route group /api/v1.
func NewRouter(appEnv string) *gin.Engine {
	if strings.ToLower(strings.TrimSpace(appEnv)) == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// Route group /api/v1
	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", func(c *gin.Context) {
			response.OK(c, "server is running", nil)
		})
	}

	return router
}
