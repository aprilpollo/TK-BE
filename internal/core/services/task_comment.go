package services

import (
	"context"

	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/core/ports/input"
	"aprilpollo/internal/core/ports/output"
	"aprilpollo/internal/pkg/query"
)

type taskCommentService struct {
	repo output.TaskCommentRepository
}

func NewTaskCommentService(repo output.TaskCommentRepository) input.TaskCommentService {
	return &taskCommentService{repo: repo}
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
