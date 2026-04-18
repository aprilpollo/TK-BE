package repository

import (
	"context"
	"time"

	"aprilpollo/internal/adapters/storage/orm/models"
	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/core/ports/output"
	"aprilpollo/internal/pkg/query"
	"aprilpollo/internal/pkg/query/gormq"
	"aprilpollo/internal/utils"

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

	if err := gormq.ApplyToGorm(base, opts).Preload("Status").Preload("Priority").Preload("Assigns.User").Find(&rows).Error; err != nil {
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

func (r *taskRepository) FindStatus(ctx context.Context, opts query.QueryOptions, project_id int64) ([]domain.TaskStatus, error) {
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

func (r *taskRepository) CreateStatus(ctx context.Context, req *domain.CreateTaskStatusReq) (*domain.TaskStatus, error) {
	var maxPosition *int

	r.db.WithContext(ctx).Model(&models.TaskStatusModel{}).
		Where(&models.TaskStatusModel{ProjectID: req.ProjectID}).
		Select("MAX(position)").
		Scan(&maxPosition)

	nextPosition := 1
	if maxPosition != nil {
		nextPosition = *maxPosition + 1
	}

	model := models.TaskStatusModel{
		ProjectID:   req.ProjectID,
		Name:        req.Name,
		Description: req.Description,
		Color:       req.Color,
		Position:    nextPosition,
	}

	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return nil, err
	}

	return model.ToDomain(), nil
}

func (r *taskRepository) UpdateStatus(ctx context.Context, req *domain.UpdateTaskStatusReq, status_id int64) (*domain.TaskStatus, error) {
	var model models.TaskStatusModel
	if err := r.db.WithContext(ctx).Where("id = ?", status_id).First(&model).Error; err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Model(&model).Updates(utils.StructToMap(req)).Error; err != nil {
		return nil, err
	}

	return model.ToDomain(), nil
}

func (r *taskRepository) DeleteStatus(ctx context.Context, status_id int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("status_id = ?", status_id).Delete(&models.TasksModel{}).Error; err != nil {
			return err
		}
		if err := tx.Where("id = ?", status_id).Delete(&models.TaskStatusModel{}).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *taskRepository) Create(ctx context.Context, req *domain.TaskReq, createBy int64) (*domain.Task, error) {
	now := time.Now()
	var model models.TasksModel

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var maxPosition *int
		if err := tx.Model(&models.TasksModel{}).
			Where(&models.TasksModel{ProjectID: req.ProjectID, StatusID: req.StatusID}).
			Select("MAX(position)").
			Scan(&maxPosition).Error; err != nil {
			return err
		}

		nextPosition := 1
		if maxPosition != nil {
			nextPosition = *maxPosition + 1
		}

		priorityID := req.PriorityID
		if priorityID == 0 {
			priorityID = 1
		}

		model = models.TasksModel{
			ProjectID:   req.ProjectID,
			Title:       req.Title,
			Description: req.Description,
			StatusID:    req.StatusID,
			PriorityID:  priorityID,
			ParentID:    req.ParentID,
			DueDate:     req.DueDate,
			Position:    nextPosition,
		}

		if err := tx.Create(&model).Error; err != nil {
			return err
		}

		if len(req.AssigneeIDs) > 0 {
			var taskAssignees []models.TaskAssignModel
			for _, assigneeID := range req.AssigneeIDs {
				taskAssignees = append(taskAssignees, models.TaskAssignModel{
					TaskID:    model.ID,
					UserID:    assigneeID,
					InvitedBy: &createBy,
					InvitedAt: &now,
					JoinedAt:  &now,
				})
			}
			if err := tx.Create(&taskAssignees).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return model.ToDomain(), nil
}

func (r *taskRepository) Update(ctx context.Context, req *domain.UpdateTaskReq, task_id int64) (*domain.Task, error) {
	var model models.TasksModel
	if err := r.db.WithContext(ctx).Where("id = ?", task_id).First(&model).Error; err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Model(&model).Updates(utils.StructToMap(req)).Error; err != nil {
		return nil, err
	}

	return model.ToDomain(), nil
}

func (r *taskRepository) Delete(ctx context.Context, task_id int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("parent_id = ?", task_id).Delete(&models.TasksModel{}).Error; err != nil {
			return err
		}

		if err := tx.Where("id = ?", task_id).Delete(&models.TasksModel{}).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *taskRepository) ReorderStatus(ctx context.Context, req *domain.ReqReorderTaskStatus, project_id int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range req.Updates {
			if err := tx.Model(&models.TaskStatusModel{}).
				Where("id = ? AND project_id = ?", item.ID, project_id).
				Update("position", item.Position).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *taskRepository) ReorderTask(ctx context.Context, req *domain.ReqReorderTask, project_id int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range req.Updates {
			var status models.TaskStatusModel
			if err := tx.Where("uuid = ? AND project_id = ?", item.StatusID, project_id).First(&status).Error; err != nil {
				return err
			}

			if err := tx.Model(&models.TasksModel{}).
				Where("id = ?", item.ID).
				Updates(map[string]interface{}{
					"position":  item.Position,
					"status_id": status.ID,
				}).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
