package output

import (
	"context"

	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/pkg/query"
)

type OrganizationRepository interface {
	FindAll(ctx context.Context, opts query.QueryOptions) ([]domain.Organization, int64, error)
	FindByID(ctx context.Context, id int64) (*domain.Organization, error)
	CreateWithOwner(ctx context.Context, org *domain.Organization, ownerUserID int64) error
	Update(ctx context.Context, id int64, req *domain.UpdateOrganizationReq) error
	Delete(ctx context.Context, id int64) error

	FindByUserID(ctx context.Context, userID int64, opts query.QueryOptions) ([]domain.UserOrganization, int64, error)
	FindByUserIDWithPrimaryDetails(ctx context.Context, userID int64) ([]domain.UserOrganizationWithDetail, error)
	FindPrimaryOrgWithDetails(ctx context.Context, userID int64) (*domain.PrimaryOrgPermissions, error)
	FindMembers(ctx context.Context, orgID int64, opts query.QueryOptions) ([]domain.OrganizationMember, int64, error)
	CreateMember(ctx context.Context, member *domain.OrganizationMember) error
	UpdateMember(ctx context.Context, orgID int64, memberID int64, req *domain.UpdateMemberReq) error
	DeleteMember(ctx context.Context, orgID int64, memberID int64) error
}
