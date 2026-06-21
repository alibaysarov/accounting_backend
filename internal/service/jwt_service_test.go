package service

import (
	"acc_backend/internal/dto"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestJwtService_GeneratesTokenPair(t *testing.T) {
	service := makeJwtService()

	userId := "user_1"
	tokenPair, err := service.GenerateTokenPair(userId)
	if err != nil {
		t.Fatalf("Test failed! An error occured %v", err)
	}
	if tokenPair == nil {
		t.Fatalf("Test failed! Token pair is empty ")
	}
}

func TestJwtService_VerifyToken(t *testing.T) {
	service := makeJwtService()
	userId := "user_1"

	tokenStr, err := service.generateAccessToken(userId)
	if err != nil {
		t.Fatalf("Test Failed! An error occured %v", err)
	}

	resultId, err := service.Verify(tokenStr)
	if err != nil {
		t.Fatalf("Test Failed! An error occured %v", err)
	}
	if resultId != userId {
		t.Fatalf("Test Failed! expected value is not equal")
	}
}

func TestJwtService_Refresh(t *testing.T) {
	service := makeJwtService()

	tests := []struct {
		name         string
		accessToken  string
		refreshToken string
		wantErr      bool
		errContains  string
	}{
		{
			name:         "Success",
			accessToken:  getAccessToken(service, "id"),
			refreshToken: getRefreshToken(service, "id"),
			wantErr:      false,
		},
		{
			name:         "Throw Expired refresh token",
			accessToken:  getAccessToken(service, "id"),
			refreshToken: getExpiredRefreshToken(service, "id"),
			wantErr:      true,
			errContains:  "token is expired",
		},
		{
			name:         "Throw Error for corrupted token",
			accessToken:  getAccessToken(service, "id"),
			refreshToken: getExpiredRefreshToken(service, "id"),
			wantErr:      true,
			errContains:  "token is expired",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dto := dto.TokenPair{
				AccessToken:  tt.accessToken,
				RefreshToken: tt.refreshToken,
			}
			result, err := service.Refresh(&dto)

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
				t.Fatal("expected user, got nil")
			}
		})
	}
}

func getCorruptedAccessToken(s *JwtService) string {
	now := time.Now()
	week := 168 * time.Hour
	exp := now.Add(-week)
	token, _ := s.getToken(&jwt.MapClaims{
		"iat": now.Unix(),
		"exp": exp.Unix(),
	})
	return token
}

func getExpiredRefreshToken(s *JwtService, userId string) string {
	now := time.Now()
	week := 168 * time.Hour
	exp := now.Add(-week)

	token, _ := s.getToken(&jwt.MapClaims{
		"user_id": userId,
		"iat":     now.Unix(),
		"exp":     exp.Unix(),
	})
	return token
}

func getAccessToken(service *JwtService, userId string) string {
	tokenStr, _ := service.generateAccessToken(userId)
	return tokenStr
}

func getRefreshToken(service *JwtService, userId string) string {
	tokenStr, _ := service.generateRefreshToken(userId)
	return tokenStr
}

func makeJwtService() *JwtService {
	service := &JwtService{
		secretKey: []byte("secret_key"),
	}

	return service
}
