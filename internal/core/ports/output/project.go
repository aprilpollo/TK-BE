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
	FindByID(ctx context.Context, projectId int64, orgId int64) (*domain.Project, error)
	FindByKey(ctx context.Context, key uuid.UUID, orgId int64) (*domain.Project, error)
	Create(ctx context.Context, orgId int64, project *domain.CreateProjectReq) (*domain.Project, error)
	CreateNotificationSettings(ctx context.Context, projectId int64) (*domain.ProjectNotificationSettings, error)
	Update(ctx context.Context, projectId int64, orgId int64, req *domain.UpdateProjectReq) error
	Delete(ctx context.Context, projectId int64, orgId int64) error
	GetNotificationSettings(ctx context.Context, projectId int64) (*domain.ProjectNotificationSettings, error)
	UpdateNotificationSettings(ctx context.Context, projectId int64, req *domain.UpdateProjectNotificationSettingsReq) error
}
