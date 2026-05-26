package services

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	"github.com/google/uuid"

	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/core/ports/input"
	"aprilpollo/internal/core/ports/output"
	"aprilpollo/internal/pkg/query"
)

type taskCommentService struct {
	repo      output.TaskCommentRepository
	minio     output.FileStorage
	publisher output.CommentPublisher
}

func NewTaskCommentService(repo output.TaskCommentRepository, minio output.FileStorage, publisher output.CommentPublisher) input.TaskCommentService {
	return &taskCommentService{repo: repo, minio: minio, publisher: publisher}
}

func (s *taskCommentService) publish(ctx context.Context, taskID int64, event domain.CommentEvent) {
	if err := s.publisher.PublishComment(ctx, taskID, event); err != nil {
		log.Printf("task_comment: publish error task=%d err=%v", taskID, err)
	}
}

func (s *taskCommentService) List(ctx context.Context, opts query.QueryOptions, taskID int64) ([]domain.TaskComment, int64, error) {
	return s.repo.FindByTaskID(ctx, opts, taskID)
}

func (s *taskCommentService) Create(ctx context.Context, req *domain.CreateTaskCommentReq, taskID int64, userID int64) (*domain.TaskComment, error) {
	comment, err := s.repo.Create(ctx, req, taskID, userID)
	if err != nil {
		return nil, err
	}
	s.publish(ctx, taskID, domain.CommentEvent{Type: domain.CommentEventCreated, TaskID: taskID, Comment: comment})
	return comment, nil
}

func (s *taskCommentService) Update(ctx context.Context, req *domain.UpdateTaskCommentReq, commentID int64) (*domain.TaskComment, error) {
	comment, err := s.repo.Update(ctx, req, commentID)
	if err != nil {
		return nil, err
	}
	s.publish(ctx, comment.TaskID, domain.CommentEvent{Type: domain.CommentEventUpdated, TaskID: comment.TaskID, Comment: comment})
	return comment, nil
}

func (s *taskCommentService) Delete(ctx context.Context, commentID int64) error {
	comment, err := s.repo.FindByID(ctx, commentID)
	if err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, commentID); err != nil {
		return err
	}
	s.publish(ctx, comment.TaskID, domain.CommentEvent{Type: domain.CommentEventDeleted, TaskID: comment.TaskID, ID: commentID})
	return nil
}

func (s *taskCommentService) UploadFiles(ctx context.Context, items []domain.UploadFileItem, taskID int64, userID int64) ([]*domain.TaskCommentFileUploadRes, error) {
	results := make([]*domain.TaskCommentFileUploadRes, 0, len(items))

	comment, err := s.repo.Create(ctx, &domain.CreateTaskCommentReq{
		Type: domain.TaskCommentTypeComment,
	}, taskID, userID)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		ext := filepath.Ext(item.Filename)
		objectName := fmt.Sprintf("tasks/comments/%d/%s%s", taskID, uuid.New().String(), ext)

		url, err := s.minio.UploadFile(ctx, objectName, item.File, item.Size, item.ContentType)
		if err != nil {
			return nil, err
		}

		result, err := s.repo.UploadFile(ctx, &domain.TaskCommentFileUpload{
			TaskCommentID: comment.ID,
			URL:           url,
			Name:          item.Filename,
			MimeType:      item.ContentType,
			Size:          item.Size,
		})
		if err != nil {
			return nil, err
		}

		results = append(results, result)
	}

	// fetch full comment with all attached files before publishing
	full, err := s.repo.FindByID(ctx, comment.ID)
	if err == nil {
		s.publish(ctx, taskID, domain.CommentEvent{Type: domain.CommentEventCreated, TaskID: taskID, Comment: full})
	}

	return results, nil
}
