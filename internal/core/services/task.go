package services

import (
	"context"
	"fmt"
	"path/filepath"

	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/core/ports/input"
	"aprilpollo/internal/core/ports/output"
	"aprilpollo/internal/pkg/query"

	"github.com/google/uuid"
)

type taskService struct {
	repo        output.TaskRepository
	commentRepo output.TaskCommentRepository
	minio       output.FileStorage
}

func NewTaskService(repo output.TaskRepository, commentRepo output.TaskCommentRepository, minio output.FileStorage) input.TaskService {
	return &taskService{repo: repo, commentRepo: commentRepo, minio: minio}
}

func (s *taskService) List(ctx context.Context, opts query.QueryOptions, projectID int64, statusID int64) ([]domain.Task, int64, error) {
	return s.repo.Find(ctx, opts, projectID, statusID)
}

func (s *taskService) GetByKey(ctx context.Context, key string) (*domain.Task, error) {
	return s.repo.FindByKey(ctx, key)
}

func (s *taskService) ListPriority(ctx context.Context) ([]domain.TaskPriority, error) {
	return s.repo.FindPriority(ctx)
}

func (s *taskService) ListStatus(ctx context.Context, opts query.QueryOptions, projectID int64) ([]domain.TaskStatus, error) {
	return s.repo.FindStatus(ctx, opts, projectID)
}

func (s *taskService) CreateStatus(ctx context.Context, req *domain.CreateTaskStatusReq) (*domain.TaskStatus, error) {
	return s.repo.CreateStatus(ctx, req)
}

func (s *taskService) CreateListStatus(ctx context.Context, projectID int64, req []domain.CreateListTaskStatusReq) error {
	return s.repo.CreateListStatus(ctx, projectID, req)
}

func (s *taskService) UpdateStatus(ctx context.Context, req *domain.UpdateTaskStatusReq, statusID int64) (*domain.TaskStatus, error) {
	return s.repo.UpdateStatus(ctx, req, statusID)
}

func (s *taskService) DeleteStatus(ctx context.Context, statusID int64) error {
	return s.repo.DeleteStatus(ctx, statusID)
}

func (s *taskService) Create(ctx context.Context, req *domain.CreateTaskReq, createdBy int64) (*domain.Task, error) {
	task, err := s.repo.Create(ctx, req, createdBy)
	if err != nil {
		return nil, err
	}

	_, _ = s.commentRepo.Create(ctx, &domain.CreateTaskCommentReq{
		Type:   domain.TaskCommentTypeEvent,
		Action: "created_task",
	}, task.ID, createdBy)

	return task, nil
}

func (s *taskService) Update(ctx context.Context, req *domain.UpdateTaskReq, taskID int64) (*domain.Task, error) {
	return s.repo.Update(ctx, req, taskID)
}

func (s *taskService) Delete(ctx context.Context, taskID int64) error {
	return s.repo.Delete(ctx, taskID)
}

func (s *taskService) ReorderStatus(ctx context.Context, req *domain.ReorderTaskStatusReq, projectID int64) error {
	return s.repo.ReorderStatus(ctx, req, projectID)
}

func (s *taskService) ReorderTask(ctx context.Context, req *domain.ReorderTaskReq, projectID int64) error {
	return s.repo.ReorderTask(ctx, req, projectID)
}

func (s *taskService) ListByToday(ctx context.Context, opts query.QueryOptions, userID int64, orgID int64) ([]domain.TaskToday, int64, error) {
	return s.repo.FindByToday(ctx, opts, userID, orgID)
}

func (s *taskService) ListOverdue(ctx context.Context, opts query.QueryOptions, userID int64, orgID int64) ([]domain.TaskToday, int64, error) {
	return s.repo.FindOverdue(ctx, opts, userID, orgID)
}

func (s *taskService) ListSubTasks(ctx context.Context, taskID int64) ([]domain.SubTask, error) {
	return s.repo.FindSubTasks(ctx, taskID)
}

func (s *taskService) CreateSubTask(ctx context.Context, req *domain.SubTaskReq, createdBy int64) (*domain.SubTask, error) {
	return s.repo.CreateSubTask(ctx, req, createdBy)
}

func (s *taskService) UpdateSubTask(ctx context.Context, req *domain.UpdateSubTaskReq, subtaskID int64) (*domain.SubTask, error) {
	return s.repo.UpdateSubTask(ctx, req, subtaskID)
}

func (s *taskService) DeleteSubTask(ctx context.Context, subtaskID int64) error {
	return s.repo.DeleteSubTask(ctx, subtaskID)
}

func (s *taskService) ReorderSubTask(ctx context.Context, req *domain.ReorderSubTaskReq, taskID int64) error {
	return s.repo.ReorderSubTask(ctx, req, taskID)
}

func (s *taskService) CreateAttachments(ctx context.Context, items []domain.UploadFileItem, taskID int64, uploadedBy int64) ([]*domain.TaskAttachmentFileUploadRes, error) {
	results := make([]*domain.TaskAttachmentFileUploadRes, 0, len(items))

	for _, item := range items {
		ext := filepath.Ext(item.Filename)
		objectName := fmt.Sprintf("tasks/attachments/%d/%s%s", taskID, uuid.New().String(), ext)

		url, err := s.minio.UploadFile(ctx, objectName, item.File, item.Size, item.ContentType)
		if err != nil {
			return nil, err
		}

		result, err := s.repo.CreateAttachments(ctx, &domain.TaskAttachment{
			TaskID:     taskID,
			Filename:   item.Filename,
			FilePath:   url,
			FileSize:   item.Size,
			MimeType:   item.ContentType,
			UploadedBy: uploadedBy,
		}, taskID)
		if err != nil {
			return nil, err
		}

		results = append(results, result)
	}

	return results, nil
}

func (s *taskService) GetAttachments(ctx context.Context, taskID int64) ([]domain.TaskAttachmentItem, error) {
	return s.repo.FindAttachments(ctx, taskID)
}

func (s *taskService) DeleteAttachment(ctx context.Context, attachmentID int64) error {
	att, err := s.repo.FindAttachment(ctx, attachmentID)
	if err != nil {
		return err
	}
	if err := s.minio.DeleteFileByURL(ctx, att.FilePath); err != nil {
		return err
	}
	return s.repo.DeleteAttachment(ctx, attachmentID)
}
