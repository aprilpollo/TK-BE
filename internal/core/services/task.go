package services

import (
	"context"

	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/core/ports/input"
	"aprilpollo/internal/core/ports/output"
	"aprilpollo/internal/pkg/query"
)

type taskService struct {
	repo output.TaskRepository
}

func NewTaskService(repo output.TaskRepository) input.TaskService {
	return &taskService{repo: repo}
}

func (s *taskService) List(ctx context.Context, opts query.QueryOptions, project_id int64, status_id int64) ([]domain.Task, int64, error) {
	return s.repo.Find(ctx, opts, project_id, status_id)
}

func (s *taskService) ListPriority(ctx context.Context) ([]domain.TaskPriority, error) {
	return s.repo.FindPriority(ctx)
}

func (s *taskService) ListStatus(ctx context.Context, project_id int64) ([]domain.TaskStatus, error) {
	return s.repo.FindStatus(ctx, project_id)
}

func (s *taskService) CreateStatus(ctx context.Context, req *domain.CreateTaskStatusReq) (*domain.TaskStatus, error) {
	return s.repo.CreateStatus(ctx, req)
}
