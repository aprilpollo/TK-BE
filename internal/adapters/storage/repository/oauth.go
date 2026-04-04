package repository

import (
	"context"
	"errors"

	"aprilpollo/internal/adapters/storage/orm/models"
	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/core/ports/output"

	"gorm.io/gorm"
)

type oauthRepository struct {
	db *gorm.DB
}

func NewOauthRepository(db *gorm.DB) output.OauthRepository {
	return &oauthRepository{db: db}
}

func (r *oauthRepository) BasicLogin(ctx context.Context, oauth *domain.BasicLogin) (*domain.OauthUser, error) {
	var row models.OauthModel
	if err := r.db.WithContext(ctx).
		Where("email = ?", oauth.Email).
		Where("provider = ?", domain.OauthProviderBasic).
		Preload("User").
		First(&row).Error; err != nil {
		return nil, err
	}
	return row.ToDomain(), nil
}

func (r *oauthRepository) UpsertSocialUser(ctx context.Context, info *domain.GoogleUserInfo, provider domain.OauthProvider) (*domain.OauthUser, error) {
	var row models.OauthModel

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// If there is oauth record → update user info and return
		err := tx.Where("provider = ? AND provider_id = ?", provider, info.ProviderID).
			Preload("User").First(&row).Error

		if err == nil {
			return tx.Model(row.User).Updates(map[string]any{
				"first_name":   info.FirstName,
				"last_name":    info.LastName,
				"display_name": info.DisplayName,
				"avatar":       info.Avatar,
			}).Error
		}

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		// Find user by email, create new one if not exists
		var user models.UserModel
		err = tx.Where("email = ?", info.Email).First(&user).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			user = models.UserModel{
				Email:       info.Email,
				FirstName:   info.FirstName,
				LastName:    info.LastName,
				DisplayName: info.DisplayName,
				Avatar:      info.Avatar,
				IsActive:    true,
				IsVerified:  true,
			}
			if err := tx.Create(&user).Error; err != nil {
				return err
			}
		}

		// Create oauth record
		row = models.OauthModel{
			UserID:     user.ID,
			Provider:   models.OauthProvider(provider),
			ProviderID: info.ProviderID,
			Email:      info.Email,
		}
		row.User = &user

		return tx.Create(&row).Error
	})

	if err != nil {
		return nil, err
	}

	return row.ToDomain(), nil
}

func (r *oauthRepository) Register(ctx context.Context) {}

func (r *oauthRepository) RefreshToken(ctx context.Context) {}
