package services

import (
	"context"
	"errors"
	"fmt"

	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/core/ports/input"
	"aprilpollo/internal/core/ports/output"
	"aprilpollo/internal/pkg/query"
	"aprilpollo/internal/utils"

	"github.com/google/uuid"
)

type userService struct {
	repo    output.UserRepository
	orgRepo output.OrganizationRepository
	minio   output.FileStorage
}

func NewUserService(repo output.UserRepository, orgRepo output.OrganizationRepository, minio output.FileStorage) input.UserService {
	return &userService{repo: repo, orgRepo: orgRepo, minio: minio}
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

func (s *userService) UpdateAvatar(ctx context.Context, userID int64, file *domain.AvatarUploadReq) error {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// remove old avatar if exists
	if user.Avatar != nil && *user.Avatar != "" {
		fmt.Printf("Deleting old avatar: %s\n", *user.Avatar)
		s.minio.DeleteFile(ctx, *user.Avatar)
	}

	if file.ContentType != "image/webp" {
		file.File, file.Size, err = utils.ConvertToWebP(file.File, 80)
		if err != nil {
			return err
		}
		file.ContentType = "image/webp"
	}

	objectName := fmt.Sprintf("avatars/%d/%s.webp", userID, uuid.New().String())

	url, err := s.minio.UploadFile(ctx, objectName, file.File, file.Size, file.ContentType)
	if err != nil {
		return err
	}

	_, err = s.Update(ctx, userID, &domain.UpdateUserReq{
		Avatar: &url,
	})
	if err != nil {
		return err
	}
	return nil
}
