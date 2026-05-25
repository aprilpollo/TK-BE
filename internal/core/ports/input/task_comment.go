package input

import (
	"context"
	"io"

	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/pkg/query"
)

type TaskCommentService interface {
	List(ctx context.Context, opts query.QueryOptions, taskID int64) ([]domain.TaskComment, int64, error)
	Create(ctx context.Context, req *domain.CreateTaskCommentReq, taskID int64, userID int64) (*domain.TaskComment, error)
	Update(ctx context.Context, req *domain.UpdateTaskCommentReq, commentID int64) (*domain.TaskComment, error)
	Delete(ctx context.Context, commentID int64) error
	UploadFile(ctx context.Context, file io.Reader, size int64, contentType string, filename string, taskID int64, userID int64) (*domain.TaskCommentFileUploadRes, error)
}
