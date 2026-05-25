package services

import (
	"context"
	"fmt"
	"io"
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

func (s *taskCommentService) UploadFile(ctx context.Context, file io.Reader, size int64, contentType string, filename string, taskID int64, userID int64) (*domain.TaskCommentFileUploadRes, error) {
	ext := filepath.Ext(filename)
	objectName := fmt.Sprintf("tasks/comments/%d/%s%s", taskID, uuid.New().String(), ext)

	url, err := s.minio.UploadFile(ctx, objectName, file, size, contentType)
	if err != nil {
		return nil, err
	}

	comment, err := s.repo.Create(ctx, &domain.CreateTaskCommentReq{
		Type: domain.TaskCommentTypeComment,
	}, taskID, userID)
	if err != nil {
		return nil, err
	}

	commentFileUpload := &domain.TaskCommentFileUpload{
		TaskCommentID: comment.ID,
		URL:           url,
		Name:          filename,
		MimeType:      contentType,
	}

	return s.repo.UploadFile(ctx, commentFileUpload)
}
