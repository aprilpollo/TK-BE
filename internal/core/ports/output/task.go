package output

import (
	"context"

	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/pkg/query"
)

type TaskRepository interface {
	Find(ctx context.Context, opts query.QueryOptions, project_id int64, status_id int64) ([]domain.Task, int64, error)
	FindPriority(ctx context.Context) ([]domain.TaskPriority, error)
	FindStatus(ctx context.Context, opts query.QueryOptions, project_id int64) ([]domain.TaskStatus, error)
	CreateStatus(ctx context.Context, req *domain.CreateTaskStatusReq) (*domain.TaskStatus, error)
	CreateListStatus(ctx context.Context, project_id int64, req []domain.CreateListTaskStatusReq) error
	UpdateStatus(ctx context.Context, req *domain.UpdateTaskStatusReq, status_id int64) (*domain.TaskStatus, error)
	DeleteStatus(ctx context.Context, status_id int64) error
	Create(ctx context.Context, req *domain.TaskReq, createBy int64) (*domain.Task, error)
	Update(ctx context.Context, req *domain.UpdateTaskReq, task_id int64) (*domain.Task, error)
	Delete(ctx context.Context, task_id int64) error
	ReorderStatus(ctx context.Context, req *domain.ReqReorderTaskStatus, project_id int64) error
	ReorderTask(ctx context.Context, req *domain.ReqReorderTask, project_id int64) error
}
