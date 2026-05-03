package services

import (
	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/core/ports/input"
	"aprilpollo/internal/core/ports/output"
	"aprilpollo/internal/pkg/query"
	"aprilpollo/internal/utils"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

type projectService struct {
	repo     output.ProjectRepository
	taskRepo output.TaskRepository
	minio    output.FileStorage
}

func NewProjectService(repo output.ProjectRepository, taskRepo output.TaskRepository, minio output.FileStorage) input.ProjectService {
	return &projectService{repo: repo, taskRepo: taskRepo, minio: minio}
}

func (s *projectService) List(ctx context.Context, opts query.QueryOptions, orgId int64) ([]domain.Project, int64, error) {
	return s.repo.FindAll(ctx, opts, orgId)
}

func (s *projectService) ListStatuses(ctx context.Context) ([]domain.ProjectStatus, error) {
	return s.repo.FindStatuses(ctx)
}

func (s *projectService) GetByID(ctx context.Context, projectId int64, orgId int64) (*domain.Project, error) {
	return s.repo.FindByID(ctx, projectId, orgId)
}

func (s *projectService) GetByKey(ctx context.Context, key uuid.UUID, orgId int64) (*domain.Project, error) {
	return s.repo.FindByKey(ctx, key, orgId)
}

func (s *projectService) Create(ctx context.Context, orgId int64, req *domain.CreateProjectReq) (*domain.Project, error) {
	pRepo, err := s.repo.Create(ctx, orgId, req)
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

func (s *projectService) Update(ctx context.Context, projectId int64, orgId int64, req *domain.UpdateProjectReq) error {
	return s.repo.Update(ctx, projectId, orgId, req)
}

func (s *projectService) UpdateLogo(ctx context.Context, projectId int64, orgId int64, file *domain.LogoUploadReq) error {
	project, err := s.repo.FindByID(ctx, projectId, orgId)
	if err != nil {
		return err
	}
	if project == nil {
		return errors.New("project not found")
	}

	if project.LogoURL != nil && *project.LogoURL != "" {
		if objectName, err := utils.ExtractObjectName(*project.LogoURL); err == nil {
			if err := s.minio.DeleteFile(ctx, objectName); err != nil {
				fmt.Printf("Failed to delete old avatar: %v\n", err)
			}
		}
	}

	if file.ContentType != "image/webp" {
		file.File, file.Size, err = utils.ConvertToWebP(file.File, 80)
		if err != nil {
			return err
		}
		file.ContentType = "image/webp"
	}

	objectName := fmt.Sprintf("project/logos/%d/%s.webp", projectId, uuid.New().String())

	url, err := s.minio.UploadFile(ctx, objectName, file.File, file.Size, file.ContentType)
	if err != nil {
		return err
	}

	err = s.Update(ctx, projectId, orgId, &domain.UpdateProjectReq{
		LogoURL: &url,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *projectService) Delete(ctx context.Context, projectId int64, orgId int64) error {
	return s.repo.Delete(ctx, projectId, orgId)
}
