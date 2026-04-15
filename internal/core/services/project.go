package services

import (
	"context"

	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/core/ports/input"
	"aprilpollo/internal/core/ports/output"
	"aprilpollo/internal/pkg/query"
	"github.com/google/uuid"
)

type projectService struct {
	repo output.ProjectRepository
}

func NewProjectService(repo output.ProjectRepository) input.ProjectService {
	return &projectService{repo: repo}
}

func (s *projectService) List(ctx context.Context, opts query.QueryOptions, orgId int64) ([]domain.Project, int64, error) {
	return s.repo.FindAll(ctx, opts, orgId)
}

func (s *projectService) ListStatuses(ctx context.Context) ([]domain.ProjectStatus, error) {
	return s.repo.FindStatuses(ctx)
}

func (s *projectService) GetByID(ctx context.Context, id int64, orgId int64) (*domain.Project, error) {
	return s.repo.FindByID(ctx, id, orgId)
}

func (s *projectService) GetByKey(ctx context.Context, key uuid.UUID, orgId int64) (*domain.Project, error) {
	return s.repo.FindByKey(ctx, key, orgId)
}

func (s *projectService) Create(ctx context.Context, req *domain.CreateProjectReq) (*domain.Project, error) {
	return s.repo.Create(ctx, req)
}

func (s *projectService) Update(ctx context.Context, id int64, req *domain.UpdateProjectReq, orgId int64) error {
	return s.repo.Update(ctx, id, req, orgId)
}

func (s *projectService) Delete(ctx context.Context, id int64, orgId int64) error {
	return s.repo.Delete(ctx, id, orgId)
}
