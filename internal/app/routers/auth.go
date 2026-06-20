package routers

import (
	"acc_backend/internal/container"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(rg *gin.RouterGroup, c *container.Container) {
	auth := rg.Group("/auth")
	{
		auth.POST("/login", c.AuthHandler.Login)
	}
}
