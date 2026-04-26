package output

import (
	"context"

	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/pkg/query"
)

type CalendarRepository interface {
	Find(ctx context.Context, opts query.QueryOptions, project_id int64) ([]domain.Task, int64, error)
	FindStatus(ctx context.Context, opts query.QueryOptions, project_id int64) ([]domain.TaskStatus, error)
	FindPriority(ctx context.Context) ([]domain.TaskPriority, error)
}