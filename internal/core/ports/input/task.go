package input

import (
	"context"

	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/pkg/query"
)

type TaskService interface {
	List(ctx context.Context, opts query.QueryOptions, project_id int64, status_id int64) ([]domain.Task, int64, error)
	ListPriority(ctx context.Context) ([]domain.TaskPriority, error)
	ListStatus(ctx context.Context, project_id int64) ([]domain.TaskStatus, error)
	CreateStatus(ctx context.Context, req *domain.CreateTaskStatusReq) (*domain.TaskStatus, error)

	ReorderStatus(ctx context.Context, req *domain.ReqReorderTaskStatus, project_id int64) error
}
