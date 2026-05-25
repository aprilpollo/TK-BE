package services

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/google/uuid"

	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/core/ports/input"
	"aprilpollo/internal/core/ports/output"
	"aprilpollo/internal/pkg/query"
)

type taskCommentService struct {
	repo  output.TaskCommentRepository
	minio output.FileStorage
}

func NewTaskCommentService(repo output.TaskCommentRepository, minio output.FileStorage) input.TaskCommentService {
	return &taskCommentService{repo: repo, minio: minio}
}

func (s *taskCommentService) List(ctx context.Context, opts query.QueryOptions, taskID int64) ([]domain.TaskComment, int64, error) {
	return s.repo.FindByTaskID(ctx, opts, taskID)
}

func (s *taskCommentService) Create(ctx context.Context, req *domain.CreateTaskCommentReq, taskID int64, userID int64) (*domain.TaskComment, error) {
	return s.repo.Create(ctx, req, taskID, userID)
}

func (s *taskCommentService) Update(ctx context.Context, req *domain.UpdateTaskCommentReq, commentID int64) (*domain.TaskComment, error) {
	return s.repo.Update(ctx, req, commentID)
}

func (s *taskCommentService) Delete(ctx context.Context, commentID int64) error {
	return s.repo.Delete(ctx, commentID)
}

func (s *taskCommentService) UploadFiles(ctx context.Context, items []domain.UploadFileItem, taskID int64, userID int64) ([]*domain.TaskCommentFileUploadRes, error) {

	results := make([]*domain.TaskCommentFileUploadRes, 0, len(items))

	comment, err := s.repo.Create(ctx, &domain.CreateTaskCommentReq{
		Type: domain.TaskCommentTypeComment,
	}, taskID, userID)
	if err != nil {
		return nil, err
	}


	for _, item := range items {
		ext := filepath.Ext(item.Filename)
		objectName := fmt.Sprintf("tasks/comments/%d/%s%s", taskID, uuid.New().String(), ext)

		url, err := s.minio.UploadFile(ctx, objectName, item.File, item.Size, item.ContentType)
		if err != nil {
			return nil, err
		}

		result, err := s.repo.UploadFile(ctx, &domain.TaskCommentFileUpload{
			TaskCommentID: comment.ID,
			URL:           url,
			Name:          item.Filename,
			MimeType:      item.ContentType,
			Size:          item.Size,
		})
		if err != nil {
			return nil, err
		}

		results = append(results, result)
	}
	return results, nil
}
