package service

import (
	"acc_backend/internal/dto"
	"acc_backend/internal/model"
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(fullName string, email string, password string) (string, error)
	FindByEmail(email string) (*model.User, error)
	GetById(ctx context.Context, id string) (*model.User, error)
}

type TokenService interface {
	GenerateTokenPair(userID string) (*dto.TokenPair, error)
	Refresh(tokenPair *dto.TokenPair) (*dto.TokenPair, error)
}

type AuthService struct {
	userRepo     UserRepository
	tokenService TokenService
}

func NewAuthService(userRepo UserRepository, tokenService TokenService) *AuthService {
	return &AuthService{userRepo: userRepo, tokenService: tokenService}
}

func (s *AuthService) Register(dto *dto.RegisterDto) (*dto.TokenPair, error) {
	user, err := s.userRepo.FindByEmail(dto.Email)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) == false {
		return nil, err
	}

	if user != nil {
		return nil, errors.New("User with this email already exists!")
	}

	hashedPass, err := hashPassword(dto.Password)
	if err != nil {
		return nil, err
	}
	userId, err := s.userRepo.Create(dto.FullName, dto.Email, hashedPass)
	if err != nil {
		return nil, err
	}
	return s.tokenService.GenerateTokenPair(userId)
}

func (s *AuthService) Login(email, password string) (*dto.TokenPair, error) {

	// логика логина
	user, err := s.userRepo.FindByEmail(email)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("User with email %s not found", email)
	}
	if comparePassword(password, user.Password) {
		pair, err := s.tokenService.GenerateTokenPair(user.ID)
		if err != nil {
			return nil, err
		}
		return pair, nil
	} else {
		return nil, errors.New("Password do no match!")
	}
}

func (s *AuthService) Refresh(tokenPair *dto.TokenPair) (*dto.TokenPair, error) {
	return nil, nil
}

func (s *AuthService) GetProfile(ctx context.Context, userId string) (*dto.ProfileDto, error) {
	user, err := s.userRepo.GetById(ctx, userId)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) == false {
		return nil, err
	}
	dto := &dto.ProfileDto{
		ID:       user.ID,
		FullName: user.FullName,
		Email:    user.Email,
	}
	return dto, nil
}

func hashPassword(password string) (string, error) {
	// GenerateFromPassword automatically generates its own secure salt
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func comparePassword(password string, hashpass string) bool {
	// GenerateFromPassword automatically generates its own secure salt
	err := bcrypt.CompareHashAndPassword([]byte(hashpass), []byte(password))
	if err != nil {
		return false
	}
	return true
}
