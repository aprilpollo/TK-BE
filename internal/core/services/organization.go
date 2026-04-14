package services

import (
	"context"
	"errors"

	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/core/ports/input"
	"aprilpollo/internal/core/ports/output"
	"aprilpollo/internal/pkg/query"
)

type organizationService struct {
	repo output.OrganizationRepository
}

func NewOrganizationService(repo output.OrganizationRepository) input.OrganizationService {
	return &organizationService{repo: repo}
}

func (s *organizationService) List(ctx context.Context, opts query.QueryOptions) ([]domain.Organization, int64, error) {
	return s.repo.FindAll(ctx, opts)
}

func (s *organizationService) GetByID(ctx context.Context, id int64) (*domain.Organization, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *organizationService) Create(ctx context.Context, req *domain.CreateOrganizationReq, ownerUserID int64) (*domain.Organization, error) {
	org := &domain.Organization{
		Name:         req.Name,
		Description:  req.Description,
		LogoURL:      req.LogoURL,
		ContactEmail: req.ContactEmail,
		IsActive:     true,
	}

	if err := s.repo.CreateWithOwner(ctx, org, ownerUserID); err != nil {
		return nil, err
	}

	return org, nil
}

func (s *organizationService) Update(ctx context.Context, id int64, req *domain.UpdateOrganizationReq) (*domain.Organization, error) {
	org, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if org == nil {
		return nil, errors.New("organization not found")
	}

	if err := s.repo.Update(ctx, id, req); err != nil {
		return nil, err
	}

	return s.repo.FindByID(ctx, id)
}

func (s *organizationService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *organizationService) ListMembers(ctx context.Context, orgID int64, opts query.QueryOptions) ([]domain.OrganizationMember, int64, error) {
	return s.repo.FindMembers(ctx, orgID, opts)
}

func (s *organizationService) InviteMember(ctx context.Context, orgID int64, req *domain.InviteMemberReq, invitedBy int64) (*domain.OrganizationMember, error) {
	member := &domain.OrganizationMember{
		OrganizationID: orgID,
		UserID:         req.UserID,
		RoleID:         req.RoleID,
		StatusID:       2, // invited
		InvitedBy:      &invitedBy,
	}

	if err := s.repo.CreateMember(ctx, member); err != nil {
		return nil, err
	}

	return member, nil
}

func (s *organizationService) UpdateMember(ctx context.Context, orgID int64, memberID int64, req *domain.UpdateMemberReq) (*domain.OrganizationMember, error) {
	members, _, err := s.repo.FindMembers(ctx, orgID, query.QueryOptions{Limit: 1, Filters: []query.Filter{
		{Field: "id", Operator: "=", Value: memberID},
	}})
	if err != nil {
		return nil, err
	}
	if len(members) == 0 {
		return nil, errors.New("member not found")
	}

	if err := s.repo.UpdateMember(ctx, orgID, memberID, req); err != nil {
		return nil, err
	}

	members[0].RoleID = req.RoleID
	members[0].StatusID = req.StatusID
	return &members[0], nil
}

func (s *organizationService) UpdatePrimary(ctx context.Context, orgID int64, memberID int64) error {
	return s.repo.UpdatePrimary(ctx, orgID, memberID)
}

func (s *organizationService) RemoveMember(ctx context.Context, orgID int64, memberID int64) error {
	return s.repo.DeleteMember(ctx, orgID, memberID)
}
