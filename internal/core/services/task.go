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

func (s *taskService) List(ctx context.Context, opts query.QueryOptions, project_id int64, status_id int64) ([]domain.Task, int64, error) {
	return s.repo.Find(ctx, opts, project_id, status_id)
}

func (s *taskService) GetByKey(ctx context.Context, key string) (*domain.Task, error) {
	return s.repo.FindByKey(ctx, key)
}

func (s *taskService) ListPriority(ctx context.Context) ([]domain.TaskPriority, error) {
	return s.repo.FindPriority(ctx)
}

func (s *taskService) ListStatus(ctx context.Context, opts query.QueryOptions, project_id int64) ([]domain.TaskStatus, error) {
	return s.repo.FindStatus(ctx, opts, project_id)
}

func (s *taskService) CreateStatus(ctx context.Context, req *domain.CreateTaskStatusReq) (*domain.TaskStatus, error) {
	return s.repo.CreateStatus(ctx, req)
}

func (s *taskService) CreateListStatus(ctx context.Context, project_id int64, req []domain.CreateListTaskStatusReq) error {
	return s.repo.CreateListStatus(ctx, project_id, req)
}

func (s *taskService) UpdateStatus(ctx context.Context, req *domain.UpdateTaskStatusReq, status_id int64) (*domain.TaskStatus, error) {
	return s.repo.UpdateStatus(ctx, req, status_id)
}

func (s *taskService) DeleteStatus(ctx context.Context, status_id int64) error {
	return s.repo.DeleteStatus(ctx, status_id)
}

func (s *taskService) Create(ctx context.Context, req *domain.TaskReq, createBy int64) (*domain.Task, error) {
	task, err := s.repo.Create(ctx, req, createBy)
	if err != nil {
		return nil, err
	}

	_, _ = s.commentRepo.Create(ctx, &domain.CreateTaskCommentReq{
		Type:   domain.TaskCommentTypeEvent,
		Action: "created_task",
	}, task.ID, createBy)

	return task, nil
}

func (s *taskService) Update(ctx context.Context, req *domain.UpdateTaskReq, task_id int64) (*domain.Task, error) {
	return s.repo.Update(ctx, req, task_id)
}

func (s *taskService) Delete(ctx context.Context, task_id int64) error {
	return s.repo.Delete(ctx, task_id)
}

func (s *taskService) ReorderStatus(ctx context.Context, req *domain.ReqReorderTaskStatus, project_id int64) error {
	return s.repo.ReorderStatus(ctx, req, project_id)
}

func (s *taskService) ReorderTask(ctx context.Context, req *domain.ReqReorderTask, project_id int64) error {
	return s.repo.ReorderTask(ctx, req, project_id)
}

func (s *taskService) ListByToday(ctx context.Context, opts query.QueryOptions, userID int64, orgID int64) ([]domain.TaskToday, int64, error) {
	return s.repo.FindByToday(ctx, opts, userID, orgID)
}

func (s *taskService) ListOverdue(ctx context.Context, opts query.QueryOptions, userID int64, orgID int64) ([]domain.TaskToday, int64, error) {
	return s.repo.FindOverdue(ctx, opts, userID, orgID)
}

func (s *taskService) ListSubTasks(ctx context.Context, taskID int64) ([]domain.SubTasks, error) {
	return s.repo.FindSubTasks(ctx, taskID)
}

func (s *taskService) CreateSubTask(ctx context.Context, req *domain.SubTaskReq, createBy int64) (*domain.SubTasks, error) {
	return s.repo.CreateSubTask(ctx, req, createBy)
}

func (s *taskService) UpdateSubTask(ctx context.Context, req *domain.UpdateSubTaskReq, subtaskID int64) (*domain.SubTasks, error) {
	return s.repo.UpdateSubTask(ctx, req, subtaskID)
}

func (s *taskService) DeleteSubTask(ctx context.Context, subtaskID int64) error {
	return s.repo.DeleteSubTask(ctx, subtaskID)
}

func (s *taskService) ReorderSubTask(ctx context.Context, req *domain.ReqReorderSubTask, taskID int64) error {
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

func (s *taskService) DownloadAttachment(ctx context.Context, attachmentID int64) (*domain.TaskAttachmentDownload, error) {
	att, err := s.repo.FindAttachment(ctx, attachmentID)
	if err != nil {
		return nil, err
	}
	reader, err := s.minio.GetFileByURL(ctx, att.FilePath)
	if err != nil {
		return nil, err
	}
	return &domain.TaskAttachmentDownload{
		Reader:   reader,
		Filename: att.Filename,
		MimeType: att.MimeType,
		Size:     att.FileSize,
	}, nil
}
