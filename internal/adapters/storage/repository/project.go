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
	"github.com/google/uuid"

	"gorm.io/gorm"
)

type projectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) output.ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) FindAll(ctx context.Context, opts query.QueryOptions, orgId int64) ([]domain.Project, int64, error) {
	var rows []models.ProjectModel
	var total int64

	base := r.db.WithContext(ctx).Model(&models.ProjectModel{}).Where("organization_id = ?", orgId)

	if err := gormq.ApplyFilters(base, opts).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := gormq.ApplyToGorm(base, opts).Preload("Status").Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	projects := make([]domain.Project, len(rows))
	for i, row := range rows {
		projects[i] = *row.ToDomain()
	}

	return projects, total, nil
}

func (r *projectRepository) FindStatuses(ctx context.Context) ([]domain.ProjectStatus, error) {
	var rows []models.ProjectStatusModel
	if err := r.db.WithContext(ctx).Find(&rows).Error; err != nil {
		return nil, err
	}

	statuses := make([]domain.ProjectStatus, len(rows))
	for i, row := range rows {
		statuses[i] = *row.ToDomain()
	}

	return statuses, nil
}

func (r *projectRepository) FindByID(ctx context.Context, projectId int64, orgId int64) (*domain.Project, error) {
	var row models.ProjectModel
	if err := r.db.WithContext(ctx).Where("id = ? AND organization_id = ?", projectId, orgId).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return row.ToDomain(), nil
}

func (r *projectRepository) FindByKey(ctx context.Context, key uuid.UUID, orgId int64) (*domain.Project, error) {
	var row models.ProjectModel
	if err := r.db.WithContext(ctx).Where("key = ? AND organization_id = ?", key, orgId).Preload("Status").First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return row.ToDomain(), nil
}

func (r *projectRepository) Create(ctx context.Context, orgId int64, project *domain.CreateProjectReq) (*domain.Project, error) {
	model := models.ProjectModel{
		OrganizationID: orgId,
		Name:           project.Name,
		Description:    project.Description,
		LogoURL:        project.LogoURL,
		StartDate:      project.StartDate,
		EndDate:        project.EndDate,
	}

	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return nil, err
	}

	return model.ToDomain(), nil
}

func (r *projectRepository) CreateNotificationSettings(ctx context.Context, projectId int64) (*domain.ProjectNotificationSettings, error) {
	settings := models.ProjectNotificationSettingModel{
		ProjectID: projectId,
	}

	if err := r.db.WithContext(ctx).Create(&settings).Error; err != nil {
		return nil, err
	}

	return settings.ToDomain(), nil
}

func (r *projectRepository) Update(ctx context.Context, projectId int64, orgId int64, req *domain.UpdateProjectReq) error {
	return r.db.WithContext(ctx).Model(&models.ProjectModel{}).Where("id = ? AND organization_id = ?", projectId, orgId).Updates(utils.StructToMap(req)).Error
}

func (r *projectRepository) Delete(ctx context.Context, projectId int64, orgId int64) error {
	return r.db.WithContext(ctx).Where("id = ? AND organization_id = ?", projectId, orgId).Delete(&models.ProjectModel{}).Error
}
