package output

import (
	"context"

	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/pkg/query"
	"github.com/google/uuid"
)

type ProjectRepository interface {
	FindAll(ctx context.Context, opts query.QueryOptions, orgId int64) ([]domain.Project, int64, error)
	FindStatuses(ctx context.Context) ([]domain.ProjectStatus, error)
	FindByID(ctx context.Context, id int64, orgId int64) (*domain.Project, error)
	FindByKey(ctx context.Context, key uuid.UUID, orgId int64) (*domain.Project, error)
	Create(ctx context.Context, project *domain.CreateProjectReq, orgId int64) (*domain.Project, error)
	Update(ctx context.Context, id int64, req *domain.UpdateProjectReq, orgId int64) error
	Delete(ctx context.Context, id int64, orgId int64) error
}