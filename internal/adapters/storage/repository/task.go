package repository

import (
	"context"
	// "errors"

	"aprilpollo/internal/adapters/storage/orm/models"
	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/core/ports/output"
	"aprilpollo/internal/pkg/query"
	"aprilpollo/internal/pkg/query/gormq"
	// "aprilpollo/internal/utils"
	// "github.com/google/uuid"

	"gorm.io/gorm"
)

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) output.TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Find(ctx context.Context, opts query.QueryOptions, project_id int64, status_id int64) ([]domain.Task, int64, error) {
	var rows []models.TasksModel
	var total int64

	base := r.db.WithContext(ctx).Model(&models.TasksModel{}).Where("project_id = ? AND status_id = ?", project_id, status_id)

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := gormq.ApplyToGorm(base, opts).Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	tasks := make([]domain.Task, len(rows))
	for i, row := range rows {
		tasks[i] = *row.ToDomain()
	}

	return tasks, total, nil
}

func (r *taskRepository) FindPriority(ctx context.Context) ([]domain.TaskPriority, error) {
	var models []models.TaskPriorityModel
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}

	domains := make([]domain.TaskPriority, len(models))
	for i, model := range models {
		domains[i] = *model.ToDomain()
	}

	return domains, nil
}

func (r *taskRepository) FindStatus(ctx context.Context, project_id int64) ([]domain.TaskStatus, error) {
	var models []models.TaskStatusModel
	if err := r.db.WithContext(ctx).Where("project_id = ?", project_id).Find(&models).Error; err != nil {
		return nil, err
	}

	domains := make([]domain.TaskStatus, len(models))
	for i, model := range models {
		domains[i] = *model.ToDomain()
	}

	return domains, nil
}

func (r *taskRepository) CreateStatus(ctx context.Context, req *domain.CreateTaskStatusReq) (*domain.TaskStatus, error) {
	model := models.TaskStatusModel{
		ProjectID:   req.ProjectID,
		Name:        req.Name,
		Description: req.Description,
		Color:       req.Color,
	}

	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return nil, err
	}

	return model.ToDomain(), nil
}
