package container

import (
	"acc_backend/internal/app/middlewares"
	"acc_backend/internal/handler"
	"acc_backend/internal/repository"
	"acc_backend/internal/service"
	"acc_backend/internal/settings"

	"gorm.io/gorm"
)

type Container struct {
	AuthHandler *handler.AuthHandler
	Utils       *Utils

	// middlewares
	AuthMiddleware *middlewares.AuthMiddleware
	// сюда же UserHandler, ProductHandler и т.д.
}

type Utils struct {
	Config *settings.AppConfig
}

func NewContainer(db *gorm.DB, utils *Utils) *Container {
	// repositories
	userRepo := repository.NewUserRepository(db)

	// services
	jwtService := service.NewJwtService(utils.Config.JwtKey)
	authService := service.NewAuthService(userRepo, jwtService)

	// handlers
	authHandler := handler.NewAuthHandler(authService)
	authMiddleware := middlewares.NewAuthMiddleware(jwtService, userRepo)

	return &Container{
		AuthHandler:    authHandler,
		Utils:          utils,
		AuthMiddleware: authMiddleware,
	}
}
