package services

import (
	"context"
	"time"

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

func (s *taskService) ListStatus(ctx context.Context, opts query.QueryOptions, project_id int64) ([]domain.TaskStatus, error) {
	return s.repo.FindStatus(ctx, opts, project_id)
}

func (s *taskService) CreateStatus(ctx context.Context, req *domain.CreateTaskStatusReq) (*domain.TaskStatus, error) {
	return s.repo.CreateStatus(ctx, req)
}

func (s *taskService) CreateListStatus(ctx context.Context, project_id int64, req []domain.CreateListTaskStatusReq) error {
	return s.repo.CreateListStatus(ctx, project_id, req)
}

func (s *taskService) UpdateStatus(ctx context.Context, req *domain.UpdateTaskStatusReq, status_id int64) (*domain.TaskStatus, error) {
	return s.repo.UpdateStatus(ctx, req, status_id)
}

func (s *taskService) DeleteStatus(ctx context.Context, status_id int64) error {
	return s.repo.DeleteStatus(ctx, status_id)
}

func (s *taskService) Create(ctx context.Context, req *domain.TaskReq, createBy int64) (*domain.Task, error) {
	return s.repo.Create(ctx, req, createBy)
}

func (s *taskService) Update(ctx context.Context, req *domain.UpdateTaskReq, task_id int64) (*domain.Task, error) {
	return s.repo.Update(ctx, req, task_id)
}

func (s *taskService) Delete(ctx context.Context, task_id int64) error {
	return s.repo.Delete(ctx, task_id)
}

func (s *taskService) ReorderStatus(ctx context.Context, req *domain.ReqReorderTaskStatus, project_id int64) error {
	return s.repo.ReorderStatus(ctx, req, project_id)
}

func (s *taskService) ReorderTask(ctx context.Context, req *domain.ReqReorderTask, project_id int64) error {
	return s.repo.ReorderTask(ctx, req, project_id)
}

func (s *taskService) ListByWeekday(ctx context.Context, opts query.QueryOptions, userID int64, orgID int64, weekday time.Weekday) ([]domain.WeekdayTask, int64, error) {
	return s.repo.FindByWeekday(ctx, opts, userID, orgID, weekday)
}
