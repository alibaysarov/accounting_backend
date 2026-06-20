package model

type User struct {
	BaseModel
	FullName string `gorm:"not null"`
	Email    string `gorm:"unique"`
	Password string
}
