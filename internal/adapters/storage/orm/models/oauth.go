package models

import (
	"time"

	"aprilpollo/internal/core/domain"
	"gorm.io/gorm"
)

type OauthProvider string

const (
	OauthProviderGoogle   OauthProvider = "google"
	OauthProviderFacebook OauthProvider = "facebook"
	OauthProviderApple    OauthProvider = "apple"
	OauthProviderBasic    OauthProvider = "basic"
)

type OauthModel struct {
	ID     int64 `gorm:"primaryKey"`
	UserID int64 `gorm:"not null;index"`

	Provider   OauthProvider `gorm:"type:varchar(20);not null;uniqueIndex:idx_provider_id,idx_provider_email"`
	ProviderID string        `gorm:"size:255;uniqueIndex:idx_provider_id"`
	Email      string        `gorm:"not null;size:255;index;uniqueIndex:idx_provider_email"`
	Password   *string       `gorm:"size:255"`

	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Relationships
	User *UserModel `gorm:"foreignKey:UserID"`
}

func (OauthModel) TableName() string {
	return "oauth"
}

func (m *OauthModel) ToDomain() *domain.OauthUser {
	return &domain.OauthUser{
		ID:          m.ID,
		Email:       m.Email,
		Password:    m.Password,
		FirstName:   m.User.FirstName,
		LastName:    m.User.LastName,
		DisplayName: m.User.DisplayName,
		Bio:         m.User.Bio,
		Avatar:      m.User.Avatar,
	}
}
