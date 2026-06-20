package service

import (
	"acc_backend/internal/dto"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtService struct {
	secretKey []byte
}

func NewJwtService(key string) *JwtService {
	return &JwtService{
		secretKey: []byte(key),
	}
}

func (s *JwtService) GenerateTokenPair(userID string) (*dto.TokenPair, error) {
	accessToken, err := s.generateAccessToken(userID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateAccessToken(userID)
	if err != nil {
		return nil, err
	}
	pair := s.makeTokenPair(accessToken, refreshToken)
	return pair, nil
}

func (s *JwtService) Refresh(tokenPair *dto.TokenPair) (*dto.TokenPair, error) {
	_, err := s.verifyToken(tokenPair.RefreshToken)
	if err != nil {
		return nil, err
	}

	accToken, err := s.verifyToken(tokenPair.AccessToken)
	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		return nil, err
	}
	userId, err := s.getUserId(accToken)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.generateAccessToken(userId)
	if err != nil {
		return nil, err
	}
	fmt.Println("get token ", accessToken)

	pair := s.makeTokenPair(accessToken, tokenPair.RefreshToken)

	return pair, err
}

func (s *JwtService) Verify(tokenString string) (string, error) {
	token, err := s.verifyToken(tokenString)
	if err != nil {
		return "", err
	}

	return s.getUserId(token)

}

func (s *JwtService) getUserId(token *jwt.Token) (string, error) {
	claims, err := s.getClaims(token)
	if err != nil {
		return "", err
	}

	if userId, exists := claims["user_id"]; exists {
		return userId.(string), nil
	}
	return "", errors.New("No user_id found in claims")
}

func (s *JwtService) getClaims(token *jwt.Token) (jwt.MapClaims, error) {
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}
	return nil, errors.New("Error during claims parse")
}

func (s *JwtService) verifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return s.secretKey, nil
	}, jwt.WithValidMethods([]string{s.getSignMethod().Alg()}))
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (s *JwtService) generateAccessToken(userID string) (string, error) {
	now := time.Now()
	return s.getToken(&jwt.MapClaims{
		"user_id": userID,
		"iat":     now.Unix(),
		"exp":     s.getAccessTokenExp(now).Unix(),
	})
}

func (s *JwtService) generateRefreshToken(userID string) (string, error) {
	now := time.Now()
	exp := s.getRefreshTokenExp(now)
	return s.getToken(&jwt.MapClaims{
		"user_id": userID,
		"iat":     now.Unix(),
		"exp":     exp.Unix(),
	})
}

func (s *JwtService) getRefreshTokenExp(now time.Time) time.Time {
	exp := now.Add(168 * time.Hour)
	return exp
}

func (s *JwtService) getAccessTokenExp(now time.Time) time.Time {
	exp := now.Add(1 * time.Hour)
	return exp
}

func (*JwtService) makeTokenPair(accessToken string, refreshToken string) *dto.TokenPair {
	pair := &dto.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return pair
}

func (s *JwtService) getToken(claims *jwt.MapClaims) (string, error) {

	var t *jwt.Token
	t = jwt.NewWithClaims(s.getSignMethod(), claims)
	token, err := t.SignedString(s.secretKey)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *JwtService) getSignMethod() jwt.SigningMethod {
	return jwt.SigningMethodHS256
}
