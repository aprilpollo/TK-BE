package repository

import (
	"context"
	"errors"

	"aprilpollo/internal/adapters/storage/orm/models"
	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/core/ports/output"
	"aprilpollo/internal/pkg/query"
	"aprilpollo/internal/pkg/query/gormq"
	"aprilpollo/internal/utils"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) output.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindAll(ctx context.Context, opts query.QueryOptions) ([]domain.User, int64, error) {
	var rows []models.UserModel
	var total int64

	base := r.db.WithContext(ctx).Model(&models.UserModel{})

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := gormq.ApplyToGorm(r.db.WithContext(ctx).Model(&models.UserModel{}), opts).Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	users := make([]domain.User, len(rows))
	for i, row := range rows {
		users[i] = *row.ToDomain()
	}

	return users, total, nil
}

func (r *userRepository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	var row models.UserModel
	if err := r.db.WithContext(ctx).First(&row, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return row.ToDomain(), nil
}

func (r *userRepository) Update(ctx context.Context, id int64, req *domain.UpdateUserReq) error {
	return r.db.WithContext(ctx).Model(&models.UserModel{}).Where("id = ?", id).
		Updates(utils.StructToMap(req)).Error
}
