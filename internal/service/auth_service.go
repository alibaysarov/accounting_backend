package service

import "acc_backend/internal/model"

type UserRepository interface {
	FindByEmail(email string) (*model.User, error)
}

type AuthService struct {
	userRepo UserRepository
}

func NewAuthService(userRepo UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Login(email, password string) (string, error) {
	// логика логина
	return "token", nil
}
