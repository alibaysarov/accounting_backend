package container

import (
	"acc_backend/internal/handler"
	"acc_backend/internal/repository"
	"acc_backend/internal/service"

	"gorm.io/gorm"
)

type Container struct {
	AuthHandler *handler.AuthHandler
	// сюда же UserHandler, ProductHandler и т.д.
}

func NewContainer(db *gorm.DB) *Container {
	// repositories
	userRepo := repository.NewUserRepository(db)

	// services
	authService := service.NewAuthService(userRepo)

	// handlers
	authHandler := handler.NewAuthHandler(authService)

	return &Container{
		AuthHandler: authHandler,
	}
}
