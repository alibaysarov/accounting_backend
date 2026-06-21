package handler

import (
	"acc_backend/internal/app/helpers"
	"acc_backend/internal/dto"
	"context"

	"github.com/gin-gonic/gin"
)

type AuthService interface {
	Register(body *dto.RegisterDto) (*dto.TokenPair, error)
	Login(email, password string) (*dto.TokenPair, error)
	GetProfile(ctx context.Context, userId string) (*dto.ProfileDto, error)
}

type AuthHandler struct {
	authService AuthService
}

func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {

	var req dto.RegisterDto
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.Fail(c, 400, err.Error())
		return
	}

	token, err := h.authService.Register(&req)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"token": token})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"token": token})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {

}
