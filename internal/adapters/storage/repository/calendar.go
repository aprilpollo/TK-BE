package repository

import (
	"context"
	// "fmt"
	// "time"

	"aprilpollo/internal/adapters/storage/orm/models"
	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/core/ports/output"
	"aprilpollo/internal/pkg/query"
	"aprilpollo/internal/pkg/query/gormq"
	// "aprilpollo/internal/utils"

	"gorm.io/gorm"
)

type calendarRepository struct {
	db *gorm.DB
}

func NewCalendarRepository(db *gorm.DB) output.CalendarRepository {
	return &calendarRepository{db: db}
}

func (r *calendarRepository) Find(ctx context.Context, opts query.QueryOptions, project_id int64) ([]domain.Task, int64, error) {
	var rows []models.TasksModel
	var total int64

	base := r.db.WithContext(ctx).Model(&models.TasksModel{}).Where("project_id = ?", project_id)

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := gormq.ApplyToGorm(base, opts).Preload("Status").Preload("Priority").Preload("Assigns.User").Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	tasks := make([]domain.Task, len(rows))
	for i, row := range rows {
		tasks[i] = *row.ToDomain()
	}

	return tasks, total, nil
}

func (r *calendarRepository) FindStatus(ctx context.Context, opts query.QueryOptions, project_id int64) ([]domain.TaskStatus, error) {
	var models []models.TaskStatusModel
	base := r.db.WithContext(ctx).Where("project_id = ?", project_id)

	if err := gormq.ApplyToGorm(base, opts).Find(&models).Error; err != nil {
		return nil, err
	}

	domains := make([]domain.TaskStatus, len(models))
	for i, model := range models {
		domains[i] = *model.ToDomain()
	}

	return domains, nil
}

func (r *calendarRepository) FindPriority(ctx context.Context) ([]domain.TaskPriority, error) {
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
