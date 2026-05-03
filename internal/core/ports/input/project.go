package input

import (
	"context"

	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/pkg/query"

	"github.com/google/uuid"
)

type ProjectService interface {
	List(ctx context.Context, opts query.QueryOptions, orgId int64) ([]domain.Project, int64, error)
	ListStatuses(ctx context.Context) ([]domain.ProjectStatus, error)
	GetByID(ctx context.Context, projectId int64, orgId int64) (*domain.Project, error)
	GetByKey(ctx context.Context, key uuid.UUID, orgId int64) (*domain.Project, error)
	Create(ctx context.Context, orgId int64, req *domain.CreateProjectReq) (*domain.Project, error)
	Update(ctx context.Context, projectId int64, orgId int64, req *domain.UpdateProjectReq) error
	UpdateLogo(ctx context.Context, projectId int64, orgId int64, file *domain.LogoUploadReq) error
	Delete(ctx context.Context, projectId int64, orgId int64) error
}
