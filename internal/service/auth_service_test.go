package service

import (
	"acc_backend/internal/dto"
	"acc_backend/internal/model"
	"context"
	"errors"
	"strings"
	"testing"

	"gorm.io/gorm"
)

type MockUserRepository struct {
	CreateFn      func(fullName, email, password string) (string, error)
	FindByEmailFn func(email string) (*model.User, error)
	GetByIdFn     func(ctx context.Context, id string) (*model.User, error)
}

func (m *MockUserRepository) Create(fullName, email, password string) (string, error) {
	return m.CreateFn(fullName, email, password)
}

func (m *MockUserRepository) FindByEmail(email string) (*model.User, error) {
	return m.FindByEmailFn(email)
}

func (m *MockUserRepository) GetById(ctx context.Context, id string) (*model.User, error) {
	return m.GetByIdFn(ctx, id)
}

func TestAuthService_Register(t *testing.T) {

	tests := []struct {
		name        string
		findByEmail func(email string) (*model.User, error)
		create      func(fullName, email, password string) (string, error)
		wantErr     bool
		errContains string
	}{
		{
			name: "success",
			findByEmail: func(email string) (*model.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
			create: func(fullName, email, password string) (string, error) {
				return "user-id-123", nil
			},
			wantErr: false,
		},
		{
			name: "user already exists",
			findByEmail: func(email string) (*model.User, error) {

				return &model.User{
					BaseModel: model.BaseModel{ID: "existing-id"},
					Email:     email}, nil
			},
			wantErr:     true,
			errContains: "already exists",
		},
		{
			name: "repo error on find",
			findByEmail: func(email string) (*model.User, error) {
				return nil, errors.New("db connection failed")
			},
			wantErr:     true,
			errContains: "db connection failed",
		},
		{
			name: "repo error on create",
			findByEmail: func(email string) (*model.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
			create: func(fullName, email, password string) (string, error) {
				return "", errors.New("insert failed")
			},
			wantErr:     true,
			errContains: "insert failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{
				FindByEmailFn: tt.findByEmail,
				CreateFn:      tt.create,
			}

			service := makeAuthService(mockRepo)

			result, err := service.Register(&dto.RegisterDto{
				FullName: "John Doe",
				Email:    "john@example.com",
				Password: "secret123",
			})

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Fatalf("expected error to contain %q, got %q", tt.errContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result == nil {
				t.Fatal("expected token pair, got nil")
			}
		})
	}
}

func makeAuthService(mockRepo *MockUserRepository) *AuthService {
	jwtService := makeJwtService()
	service := NewAuthService(mockRepo, jwtService)
	return service
}
