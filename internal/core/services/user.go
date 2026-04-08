package services

import (
	"context"
	"errors"

	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/core/ports/input"
	"aprilpollo/internal/core/ports/output"
	"aprilpollo/internal/pkg/query"
)

type userService struct {
	repo    output.UserRepository
	orgRepo output.OrganizationRepository
}

func NewUserService(repo output.UserRepository, orgRepo output.OrganizationRepository) input.UserService {
	return &userService{repo: repo, orgRepo: orgRepo}
}

func (s *userService) List(ctx context.Context, opts query.QueryOptions) ([]domain.User, int64, error) {
	return s.repo.FindAll(ctx, opts)
}

func (s *userService) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *userService) ListMyOrganizations(ctx context.Context, userID int64, opts query.QueryOptions) ([]domain.UserOrganization, int64, error) {
	return s.orgRepo.FindByUserID(ctx, userID, opts)
}

func (s *userService) GetMyPrimaryOrgPermissions(ctx context.Context, userID int64) (*domain.PrimaryOrgPermissions, error) {
	return s.orgRepo.FindPrimaryOrgWithDetails(ctx, userID)
}

func (s *userService) Update(ctx context.Context, id int64, req *domain.UpdateUserReq) (*domain.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	if err := s.repo.Update(ctx, id, req); err != nil {
		return nil, err
	}

	return s.repo.FindByID(ctx, id)
}
