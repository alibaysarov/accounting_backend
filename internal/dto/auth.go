package dto

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type RegisterDto struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
