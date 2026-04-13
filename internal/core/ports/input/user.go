package input

import (
	"context"

	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/pkg/query"
)

type UserService interface {
	List(ctx context.Context, opts query.QueryOptions) ([]domain.User, int64, error)
	GetByID(ctx context.Context, id int64) (*domain.User, error)
	Update(ctx context.Context, id int64, req *domain.UpdateUserReq) (*domain.User, error)
	ListMyOrganizations(ctx context.Context, userID int64, opts query.QueryOptions) ([]domain.UserOrganization, int64, error)
	GetMyPrimaryOrgPermissions(ctx context.Context, userID int64) (*domain.PrimaryOrgPermissions, error)
	UpdateAvatar(ctx context.Context, userID int64, file *domain.AvatarUploadReq) error
}
