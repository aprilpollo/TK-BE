package output

import (
	"context"

	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/pkg/query"
)

type TaskCommentRepository interface {
	FindByTaskID(ctx context.Context, opts query.QueryOptions, taskID int64) ([]domain.TaskComment, int64, error)
	FindByID(ctx context.Context, commentID int64) (*domain.TaskComment, error)
	Create(ctx context.Context, req *domain.CreateTaskCommentReq, taskID int64, userID int64) (*domain.TaskComment, error)
	Update(ctx context.Context, req *domain.UpdateTaskCommentReq, commentID int64) (*domain.TaskComment, error)
	Delete(ctx context.Context, commentID int64) error
	UploadFile(ctx context.Context, req *domain.TaskCommentFileUpload) (*domain.TaskCommentFileUploadRes, error)
}
