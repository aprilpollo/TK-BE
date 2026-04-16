package services

import (
	"context"

	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/core/ports/input"
	"aprilpollo/internal/core/ports/output"
	// "aprilpollo/internal/pkg/query"
	// "github.com/google/uuid"
)

type taskService struct {
	repo output.TaskRepository
}

func NewTaskService(repo output.TaskRepository) input.TaskService {
	return &taskService{repo: repo}
}

func (s *taskService) ListPriority(ctx context.Context) ([]domain.TaskPriority, error) {
	return s.repo.FindPriority(ctx)
}

func (s *taskService) CreateStatus(ctx context.Context, req *domain.CreateTaskStatusReq) (*domain.TaskStatus, error) {
	return s.repo.CreateStatus(ctx, req)
}
