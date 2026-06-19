package routers

import (
	"acc_backend/internal/container"

	"github.com/gin-gonic/gin"
)

func NewRouter(c *container.Container) *gin.Engine {

	router := gin.Default()
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"Message": "Pong"})
	})

	api := router.Group("/api/v1")
	{
		RegisterAuthRoutes(api)
	}

	return router
}
