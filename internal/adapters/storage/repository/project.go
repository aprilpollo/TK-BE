package repository

import (
	"context"
	"errors"
	"sort"
	"time"

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

func (r *projectRepository) GetNotificationSettings(ctx context.Context, projectId int64) (*domain.ProjectNotificationSettings, error) {
	var settings models.ProjectNotificationSettingModel
	if err := r.db.WithContext(ctx).Where("project_id = ?", projectId).First(&settings).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return settings.ToDomain(), nil
}

func (r *projectRepository) UpdateNotificationSettings(ctx context.Context, projectId int64, req *domain.UpdateProjectNotificationSettingsReq) error {
	return r.db.WithContext(ctx).Model(&models.ProjectNotificationSettingModel{}).Where("project_id = ?", projectId).Updates(utils.StructToMap(req)).Error
}

func weekStart() time.Time {
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	d := now.AddDate(0, 0, -(weekday - 1))
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

func (r *projectRepository) GetTaskSummary(ctx context.Context, projectId int64) (*domain.TaskSummary, error) {
	monday := weekStart()
	db := r.db.WithContext(ctx)

	var total, completed, totalThisWeek, completedThisWeek, pendingThisWeek int64

	if err := db.Model(&models.TasksModel{}).Where("project_id = ?", projectId).Count(&total).Error; err != nil {
		return nil, err
	}

	if err := db.Model(&models.TasksModel{}).
		Joins("JOIN task_statuses ts ON tasks.status_id = ts.id AND ts.is_complete = true AND ts.deleted_at IS NULL").
		Where("tasks.project_id = ? AND tasks.deleted_at IS NULL", projectId).
		Count(&completed).Error; err != nil {
		return nil, err
	}

	if err := db.Model(&models.TasksModel{}).
		Where("project_id = ? AND created_at >= ?", projectId, monday).
		Count(&totalThisWeek).Error; err != nil {
		return nil, err
	}

	if err := db.Model(&models.TasksModel{}).
		Joins("JOIN task_statuses ts ON tasks.status_id = ts.id AND ts.is_complete = true AND ts.deleted_at IS NULL").
		Where("tasks.project_id = ? AND tasks.deleted_at IS NULL AND tasks.updated_at >= ?", projectId, monday).
		Count(&completedThisWeek).Error; err != nil {
		return nil, err
	}

	if err := db.Model(&models.TasksModel{}).
		Joins("JOIN task_statuses ts ON tasks.status_id = ts.id AND ts.is_complete = false AND ts.deleted_at IS NULL").
		Where("tasks.project_id = ? AND tasks.deleted_at IS NULL AND tasks.updated_at >= ?", projectId, monday).
		Count(&pendingThisWeek).Error; err != nil {
		return nil, err
	}

	return &domain.TaskSummary{
		Total:             total,
		TotalThisWeek:     totalThisWeek,
		Completed:         completed,
		CompletedThisWeek: completedThisWeek,
		Pending:           total - completed,
		PendingThisWeek:   pendingThisWeek,
		Cancelled:         0,
		CancelledThisWeek: 0,
	}, nil
}

func (r *projectRepository) GetTaskVelocityChart(ctx context.Context, projectId int64) ([]domain.TaskVelocityPoint, error) {
	type dailyCount struct {
		Date  time.Time
		Count int64
	}

	var createdRows, completedRows []dailyCount

	if err := r.db.WithContext(ctx).Raw(`
		SELECT DATE(created_at) AS date, COUNT(*) AS count
		FROM tasks
		WHERE project_id = ? AND deleted_at IS NULL
		GROUP BY DATE(created_at)
		ORDER BY date
	`, projectId).Scan(&createdRows).Error; err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Raw(`
		SELECT DATE(t.updated_at) AS date, COUNT(*) AS count
		FROM tasks t
		JOIN task_statuses ts ON t.status_id = ts.id AND ts.is_complete = true AND ts.deleted_at IS NULL
		WHERE t.project_id = ? AND t.deleted_at IS NULL
		GROUP BY DATE(t.updated_at)
		ORDER BY date
	`, projectId).Scan(&completedRows).Error; err != nil {
		return nil, err
	}

	points := make(map[string]*domain.TaskVelocityPoint)
	for _, row := range createdRows {
		key := row.Date.Format("2006-01-02")
		if points[key] == nil {
			points[key] = &domain.TaskVelocityPoint{Date: key}
		}
		points[key].Created = row.Count
	}
	for _, row := range completedRows {
		key := row.Date.Format("2006-01-02")
		if points[key] == nil {
			points[key] = &domain.TaskVelocityPoint{Date: key}
		}
		points[key].Completed = row.Count
	}

	result := make([]domain.TaskVelocityPoint, 0, len(points))
	for _, p := range points {
		result = append(result, *p)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Date < result[j].Date })

	return result, nil
}

func (r *projectRepository) GetProjectMembers(ctx context.Context, projectId int64, opts query.QueryOptions) ([]domain.TaskMember, int64, error) {
	type memberRow struct {
		ID          int64
		DisplayName string
		Email       string
		Avatar      *string
	}

	const baseSQL = ` FROM task_assignments ta
		JOIN tasks t ON ta.task_id = t.id AND t.deleted_at IS NULL
		JOIN users u ON ta.user_id = u.id AND u.deleted_at IS NULL
		WHERE t.project_id = ? AND ta.deleted_at IS NULL`

	var total int64
	if err := r.db.WithContext(ctx).Raw("SELECT COUNT(DISTINCT u.id)"+baseSQL, projectId).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	limit := opts.Limit
	if limit <= 0 {
		limit = 10
	}
	offset := opts.Offset

	var rows []memberRow
	if err := r.db.WithContext(ctx).Raw(
		"SELECT DISTINCT u.id, u.display_name, u.email, u.avatar"+baseSQL+" ORDER BY u.id LIMIT ? OFFSET ?",
		projectId, limit, offset,
	).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	members := make([]domain.TaskMember, len(rows))
	for i, row := range rows {
		members[i] = domain.TaskMember{
			ID:     row.ID,
			Name:   row.DisplayName,
			Email:  row.Email,
			Avatar: row.Avatar,
		}
	}
	return members, total, nil
}

func (r *projectRepository) GetUpcomingDeadlines(ctx context.Context, projectId int64, opts query.QueryOptions) ([]domain.TaskDeadlineItem, int64, error) {
	monday := weekStart()
	sundayEnd := monday.AddDate(0, 0, 7).Add(-time.Nanosecond)
	startMs := monday.UnixMilli()
	endMs := sundayEnd.UnixMilli()

	base := r.db.WithContext(ctx).Model(&models.TasksModel{}).
		Joins("JOIN task_statuses ts ON tasks.status_id = ts.id AND ts.is_complete = false AND ts.deleted_at IS NULL").
		Where("tasks.project_id = ? AND tasks.end_date BETWEEN ? AND ? AND tasks.deleted_at IS NULL", projectId, startMs, endMs)

	var total int64
	if err := r.db.WithContext(ctx).Model(&models.TasksModel{}).
		Joins("JOIN task_statuses ts ON tasks.status_id = ts.id AND ts.is_complete = false AND ts.deleted_at IS NULL").
		Where("tasks.project_id = ? AND tasks.end_date BETWEEN ? AND ? AND tasks.deleted_at IS NULL", projectId, startMs, endMs).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	limit := opts.Limit
	if limit <= 0 {
		limit = 10
	}

	var rows []models.TasksModel
	if err := base.Preload("Status").Preload("Priority").
		Limit(limit).Offset(opts.Offset).
		Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	items := make([]domain.TaskDeadlineItem, 0, len(rows))
	for _, row := range rows {
		if row.EndDate == nil {
			continue
		}
		dueDate := time.UnixMilli(*row.EndDate).Format("2006-01-02")

		status := domain.TaskStatus{}
		if row.Status != nil {
			status = *row.Status.ToDomain()
		}
		priority := domain.TaskPriority{}
		if row.Priority != nil {
			priority = *row.Priority.ToDomain()
		}

		items = append(items, domain.TaskDeadlineItem{
			ID:       row.ID,
			Key:      row.Key.String(),
			Name:     row.Title,
			DueDate:  dueDate,
			Priority: priority,
			Status:   status,
		})
	}
	return items, total, nil
}
