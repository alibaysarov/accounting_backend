package dto

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RegisterDto struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ProfileDto struct {
	ID       string `json:"id"`
	FullName string `gorm:"not null"`
	Email    string `gorm:"unique"`
}
