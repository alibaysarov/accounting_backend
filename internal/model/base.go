package model

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	gorm.Model
	ID string `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
