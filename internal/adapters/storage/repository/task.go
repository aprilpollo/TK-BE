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

	if err := gormq.ApplyFilters(base, opts).Count(&total).Error; err != nil {
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
	base := r.db.WithContext(ctx).Where("project_id = ?", project_id).Order("is_complete ASC, position ASC")

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
	var result models.TaskStatusModel

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var maxPosition *int

		tx.Model(&models.TaskStatusModel{}).
			Where("project_id = ? AND is_complete = ?", req.ProjectID, false).
			Select("MAX(position)").
			Scan(&maxPosition)

		nextPosition := 1
		if maxPosition != nil {
			nextPosition = *maxPosition + 1
		}

		if !req.IsComplete {

			if err := tx.Model(&models.TaskStatusModel{}).
				Where("project_id = ? AND is_complete = ?", req.ProjectID, true).
				UpdateColumn("position", gorm.Expr("position + 1")).Error; err != nil {
				return err
			}
		} else {

			tx.Model(&models.TaskStatusModel{}).
				Where("project_id = ?", req.ProjectID).
				Select("MAX(position)").
				Scan(&maxPosition)

			nextPosition = 1
			if maxPosition != nil {
				nextPosition = *maxPosition + 1
			}
		}

		result = models.TaskStatusModel{
			ProjectID:   req.ProjectID,
			Name:        req.Name,
			Description: req.Description,
			Color:       req.Color,
			Position:    nextPosition,
			IsComplete:  req.IsComplete,
		}

		return tx.Create(&result).Error
	})
	if err != nil {
		return nil, err
	}

	return result.ToDomain(), nil
}

func (r *taskRepository) CreateListStatus(ctx context.Context, project_id int64, req []domain.CreateListTaskStatusReq) error {
	var toCreate []models.TaskStatusModel

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i, item := range req {
			if item.UUID != nil {
				if err := tx.Model(&models.TaskStatusModel{}).
					Where("uuid = ? AND project_id = ?", item.UUID, project_id).
					Updates(map[string]any{
						"name":        item.Name,
						"description": item.Description,
						"color":       item.Color,
						"position":    i + 1,
						"is_complete": item.IsComplete,
					}).Error; err != nil {
					return err
				}
			} else {
				toCreate = append(toCreate, models.TaskStatusModel{
					ProjectID:   project_id,
					Name:        item.Name,
					Description: item.Description,
					Color:       item.Color,
					IsComplete:  item.IsComplete,
					Position:    i + 1,
				})
			}
		}

		if len(toCreate) > 0 {
			if err := tx.Create(&toCreate).Error; err != nil {
				return err
			}
		}

		return nil
	})
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
			StartDate:   req.StartDate,
			EndDate:     req.EndDate,
			AllDay:      req.AllDay,
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

func (r *taskRepository) FindByToday(ctx context.Context, opts query.QueryOptions, userID int64, orgID int64) ([]domain.TaskToday, int64, error) {
	now := time.Now()
	unix := now.UnixMilli()

	base := r.db.WithContext(ctx).Model(&models.TasksModel{}).
		Joins("JOIN task_assignments ta ON ta.task_id = tasks.id AND ta.user_id = ? AND ta.deleted_at IS NULL", userID).
		Joins("JOIN projects p ON p.id = tasks.project_id AND p.organization_id = ? AND p.deleted_at IS NULL", orgID).
		Joins("JOIN task_statuses ts ON ts.id = tasks.status_id AND ts.deleted_at IS NULL AND ts.is_complete != ?", true).
		Where("tasks.deleted_at IS NULL AND tasks.start_date <= ? AND tasks.end_date >= ?", unix, unix)

	var total int64
	if err := gormq.ApplyFilters(base, opts).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []models.TasksModel
	if err := gormq.ApplyToGorm(base, opts).
		Preload("Status").Preload("Priority").Preload("Assigns.User").Preload("Project").
		Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	items := make([]domain.TaskToday, 0, len(rows))
	for _, row := range rows {
		status := domain.TaskStatus{ID: row.StatusID}
		if row.Status != nil {
			status = *row.Status.ToDomain()
		}
		priority := domain.TaskPriority{ID: row.PriorityID}
		if row.Priority != nil {
			priority = *row.Priority.ToDomain()
		}
		assigns := make([]domain.TaskAssign, 0, len(row.Assigns))
		for _, a := range row.Assigns {
			item := domain.TaskAssign{ID: a.UserID}
			if a.User != nil {
				item.Name = a.User.DisplayName
				item.Email = a.User.Email
				if a.User.Avatar != nil {
					item.Avatar = *a.User.Avatar
				}
			}
			assigns = append(assigns, item)
		}
		projectName := ""
		if row.Project != nil {
			projectName = row.Project.Name
		}
		items = append(items, domain.TaskToday{
			ID:          row.ID,
			Key:         row.Key,
			Title:       row.Title,
			StartDate:   row.StartDate,
			EndDate:     row.EndDate,
			AllDay:      row.AllDay != nil && *row.AllDay,
			Priority:    priority,
			Status:      status,
			Assignees:   assigns,
			ProjectID:   row.ProjectID,
			ProjectName: projectName,
		})
	}
	return items, total, nil
}

func (r *taskRepository) FindOverdue(ctx context.Context, opts query.QueryOptions, userID int64, orgID int64) ([]domain.TaskToday, int64, error) {
	now := time.Now()
	unix := now.UnixMilli()

	base := r.db.WithContext(ctx).Model(&models.TasksModel{}).
		Joins("JOIN task_assignments ta ON ta.task_id = tasks.id AND ta.user_id = ? AND ta.deleted_at IS NULL", userID).
		Joins("JOIN projects p ON p.id = tasks.project_id AND p.organization_id = ? AND p.deleted_at IS NULL", orgID).
		Joins("JOIN task_statuses ts ON ts.id = tasks.status_id AND ts.deleted_at IS NULL AND ts.is_complete != ?", true).
		Where("tasks.deleted_at IS NULL AND tasks.end_date < ?", unix)

	var total int64
	if err := gormq.ApplyFilters(base, opts).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []models.TasksModel
	if err := gormq.ApplyToGorm(base, opts).
		Preload("Status").Preload("Priority").Preload("Assigns.User").Preload("Project").
		Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	items := make([]domain.TaskToday, 0, len(rows))
	for _, row := range rows {
		status := domain.TaskStatus{ID: row.StatusID}
		if row.Status != nil {
			status = *row.Status.ToDomain()
		}
		priority := domain.TaskPriority{ID: row.PriorityID}
		if row.Priority != nil {
			priority = *row.Priority.ToDomain()
		}
		assigns := make([]domain.TaskAssign, 0, len(row.Assigns))
		for _, a := range row.Assigns {
			item := domain.TaskAssign{ID: a.UserID}
			if a.User != nil {
				item.Name = a.User.DisplayName
				item.Email = a.User.Email
				if a.User.Avatar != nil {
					item.Avatar = *a.User.Avatar
				}
			}
			assigns = append(assigns, item)
		}
		projectName := ""
		if row.Project != nil {
			projectName = row.Project.Name
		}
		items = append(items, domain.TaskToday{
			ID:          row.ID,
			Key:         row.Key,
			Title:       row.Title,
			StartDate:   row.StartDate,
			EndDate:     row.EndDate,
			AllDay:      row.AllDay != nil && *row.AllDay,
			Priority:    priority,
			Status:      status,
			Assignees:   assigns,
			ProjectID:   row.ProjectID,
			ProjectName: projectName,
		})
	}
	return items, total, nil
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
				Updates(map[string]any{
					"position":  item.Position,
					"status_id": status.ID,
				}).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
