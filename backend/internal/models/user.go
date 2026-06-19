package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	FullName     string     `gorm:"type:varchar(255);not null" json:"fullName"`
	Email        string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash string     `gorm:"type:varchar(255);not null" json:"-"`
	Role         string     `gorm:"type:varchar(50);default:'user';not null" json:"role"`
	PhoneNumber  string     `gorm:"type:varchar(50)" json:"phoneNumber"`
	CreatedAt    time.Time  `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt    time.Time  `gorm:"not null;default:now()" json:"updatedAt"`
	DeletedAt    *time.Time `gorm:"index" json:"-"`
}

type RegisterRequest struct {
	FullName    string `json:"fullName" binding:"required,min=2"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	PhoneNumber string `json:"phoneNumber"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
	User      User      `json:"user"`
}

type UpdateProfileRequest struct {
	FullName    string `json:"fullName" binding:"required,min=2"`
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
}
