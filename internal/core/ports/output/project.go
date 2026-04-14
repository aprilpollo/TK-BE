package output

import (
	"context"

	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/pkg/query"
)

type ProjectRepository interface {
	FindAll(ctx context.Context, opts query.QueryOptions, orgId int64) ([]domain.Project, int64, error)
	FindByID(ctx context.Context, id int64, orgId int64) (*domain.Project, error)
	Create(ctx context.Context, project *domain.CreateProjectReq) (*domain.Project, error)
	Update(ctx context.Context, id int64, req *domain.UpdateProjectReq, orgId int64) error
	Delete(ctx context.Context, id int64, orgId int64) error
}