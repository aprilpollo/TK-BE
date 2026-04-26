package input

import (
	"context"

	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/pkg/query"
)

type CalendarService interface {
	List(ctx context.Context, opts query.QueryOptions, project_id int64) ([]domain.Task, int64, error)
	ListStatus(ctx context.Context, opts query.QueryOptions, project_id int64) ([]domain.TaskStatus, error)
	ListPriority(ctx context.Context) ([]domain.TaskPriority, error)
}
