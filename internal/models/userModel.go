package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"-"`
	UserID    string         `gorm:"type:uuid;default:gen_random_uuid();uniqueIndex" json:"user_id"`
	FirstName string         `gorm:"size:100" json:"first_name"`
	LastName  string         `gorm:"size:100" json:"last_name"`
	Email     string         `gorm:"size:255;uniqueIndex;not null" json:"email" validate:"required,email"`
	Password  string         `json:"-" validate:"required,min=8"` //hide password
	UserType  string         `gorm:"size:20" json:"user_type" validate:"required,oneof=ADMIN USER"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
