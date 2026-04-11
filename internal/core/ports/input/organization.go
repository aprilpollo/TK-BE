package input

import (
	"context"

	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/pkg/query"
)

type OrganizationService interface {
	List(ctx context.Context, opts query.QueryOptions) ([]domain.Organization, int64, error)
	GetByID(ctx context.Context, id int64) (*domain.Organization, error)
	Create(ctx context.Context, req *domain.CreateOrganizationReq, ownerUserID int64) (*domain.Organization, error)
	Update(ctx context.Context, id int64, req *domain.UpdateOrganizationReq) (*domain.Organization, error)
	Delete(ctx context.Context, id int64) error

	ListMembers(ctx context.Context, orgID int64, opts query.QueryOptions) ([]domain.OrganizationMember, int64, error)
	InviteMember(ctx context.Context, orgID int64, req *domain.InviteMemberReq, invitedBy int64) (*domain.OrganizationMember, error)
	UpdateMember(ctx context.Context, orgID int64, memberID int64, req *domain.UpdateMemberReq) (*domain.OrganizationMember, error)
	UpdatePrimary(ctx context.Context, orgID int64, memberID int64) error
	RemoveMember(ctx context.Context, orgID int64, memberID int64) error
}
