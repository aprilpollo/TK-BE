package models

import (
	"time"

	"aprilpollo/internal/core/domain"
	"gorm.io/gorm"
)

type UserModel struct {
	ID          int64   `gorm:"primaryKey"`
	Email       string  `gorm:"not null;uniqueIndex;size:255"`
	FirstName   string  `gorm:"size:255"`
	LastName    string  `gorm:"size:255"`
	DisplayName string  `gorm:"size:255"`
	Bio         *string `gorm:"type:text"`
	Avatar      *string `gorm:"size:255"`
	IsActive    bool    `gorm:"default:true;index"`
	IsVerified  bool    `gorm:"default:false"`

	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Relationships
	Oauths []OauthModel `gorm:"foreignKey:UserID"`
}

func (UserModel) TableName() string {
	return "users"
}

func (m *UserModel) ToDomain() *domain.User {
	return &domain.User{
		ID:          m.ID,
		Email:       m.Email,
		FirstName:   m.FirstName,
		LastName:    m.LastName,
		DisplayName: m.DisplayName,
		Bio:         m.Bio,
		Avatar:      m.Avatar,
		Role:        "user",
		IsActive:    m.IsActive,
		IsVerified:  m.IsVerified,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}