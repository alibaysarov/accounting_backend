package service

import "testing"

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
	tokenStr, err := getAccessToken(service, userId)
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

func getAccessToken(service *JwtService, userId string) (string, error) {
	return service.generateAccessToken(userId)
}

func makeJwtService() *JwtService {
	service := &JwtService{
		secretKey: []byte("secret_key"),
	}

	return service
}
