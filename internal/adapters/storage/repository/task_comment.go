package repository

import (
	"context"
	"time"

	"aprilpollo/internal/adapters/storage/orm/models"
	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/core/ports/output"
	"aprilpollo/internal/pkg/query"
	"aprilpollo/internal/pkg/query/gormq"

	"gorm.io/gorm"
)

type taskCommentRepository struct {
	db *gorm.DB
}

func NewTaskCommentRepository(db *gorm.DB) output.TaskCommentRepository {
	return &taskCommentRepository{db: db}
}

func (r *taskCommentRepository) FindByTaskID(ctx context.Context, opts query.QueryOptions, taskID int64) ([]domain.TaskComment, int64, error) {
	var rows []models.TaskCommentModel
	var total int64

	base := r.db.WithContext(ctx).Model(&models.TaskCommentModel{}).Where("task_id = ?", taskID)

	if err := gormq.ApplyFilters(base, opts).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := gormq.ApplyToGorm(base, opts).Preload("User").Preload("Files").Order("created_at ASC").Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	comments := make([]domain.TaskComment, len(rows))
	for i, row := range rows {
		comments[i] = *row.ToDomain()
	}

	return comments, total, nil
}

func (r *taskCommentRepository) Create(ctx context.Context, req *domain.CreateTaskCommentReq, taskID int64, userID int64) (*domain.TaskComment, error) {
	var comment models.TaskCommentModel

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		comment = models.TaskCommentModel{
			TaskID:    taskID,
			UserID:    userID,
			Type:      models.TaskCommentModelType(req.Type),
			Text:      req.Text,
			Action:    req.Action,
			TimeStamp: time.Now().Unix(),
		}
		if err := tx.Create(&comment).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	var result models.TaskCommentModel
	if err := r.db.WithContext(ctx).Preload("User").Preload("Files").First(&result, comment.ID).Error; err != nil {
		return nil, err
	}

	return result.ToDomain(), nil
}

func (r *taskCommentRepository) UploadFile(ctx context.Context, req *domain.TaskCommentFileUpload) (*domain.TaskCommentFileUploadRes, error) {
	commentFile := models.TaskCommentFileModel{
		TaskCommentID: req.TaskCommentID,
		Url:           req.URL,
		Name:          req.Name,
		MimeType:      req.MimeType,
		Size:          req.Size,
	}

	if err := r.db.WithContext(ctx).Create(&commentFile).Error; err != nil {
		return nil, err
	}

	return &domain.TaskCommentFileUploadRes{
		URL:      commentFile.Url,
		Name:     commentFile.Name,
		MimeType: commentFile.MimeType,
		Size:     commentFile.Size,
	}, nil
}

func (r *taskCommentRepository) FindByID(ctx context.Context, commentID int64) (*domain.TaskComment, error) {
	var result models.TaskCommentModel
	if err := r.db.WithContext(ctx).Preload("User").Preload("Files").First(&result, commentID).Error; err != nil {
		return nil, err
	}
	return result.ToDomain(), nil
}

func (r *taskCommentRepository) Update(ctx context.Context, req *domain.UpdateTaskCommentReq, commentID int64) (*domain.TaskComment, error) {
	if err := r.db.WithContext(ctx).Model(&models.TaskCommentModel{}).Where("id = ?", commentID).Update("text", req.Text).Error; err != nil {
		return nil, err
	}

	var result models.TaskCommentModel
	if err := r.db.WithContext(ctx).Preload("User").Preload("Files").First(&result, commentID).Error; err != nil {
		return nil, err
	}

	return result.ToDomain(), nil
}

func (r *taskCommentRepository) Delete(ctx context.Context, commentID int64) error {
	return r.db.WithContext(ctx).Delete(&models.TaskCommentModel{}, commentID).Error
}
