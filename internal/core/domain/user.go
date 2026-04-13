package domain

import (
	"io"
	"time"
)

type User struct {
	ID          int64     `json:"id"`
	Email       string    `json:"email"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	DisplayName string    `json:"display_name"`
	Bio         *string   `json:"bio"`
	Avatar      *string   `json:"avatar"`
	IsActive    bool      `json:"is_active"`
	IsVerified  bool      `json:"is_verified"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UpdateUserReq struct {
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	DisplayName string  `json:"display_name"`
	Bio         *string `json:"bio"`
	Avatar      *string `json:"avatar"`
}

type AvatarUploadReq struct {
	File        io.Reader
	Size        int64
	ContentType string
	Filename    string
}
