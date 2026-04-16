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
	repo     output.ProjectRepository
	taskRepo output.TaskRepository
}

func NewProjectService(repo output.ProjectRepository, taskRepo output.TaskRepository) input.ProjectService {
	return &projectService{repo: repo, taskRepo: taskRepo}
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

func (s *projectService) Create(ctx context.Context, req *domain.CreateProjectReq, orgId int64) (*domain.Project, error) {
	pRepo, err := s.repo.Create(ctx, req, orgId)
	if err != nil {
		return nil, err
	}

	// Create default task statuses for the new project
	defaultStatuses := domain.CreateTaskStatusReq{
		ProjectID:   pRepo.ID,
		Name:        "To Do",
		Description: "Tasks that need to be done",
		Color:       "#52525B",
	}

	_, err = s.taskRepo.CreateStatus(ctx, &defaultStatuses)
	if err != nil {
		return nil, err
	}

	return pRepo, nil
}

func (s *projectService) Update(ctx context.Context, id int64, req *domain.UpdateProjectReq, orgId int64) error {
	return s.repo.Update(ctx, id, req, orgId)
}

func (s *projectService) Delete(ctx context.Context, id int64, orgId int64) error {
	return s.repo.Delete(ctx, id, orgId)
}
