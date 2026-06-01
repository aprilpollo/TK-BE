package output

import (
	"context"

	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/pkg/query"
)

type TaskRepository interface {
	Find(ctx context.Context, opts query.QueryOptions, projectID int64, statusID int64) ([]domain.Task, int64, error)
	FindByKey(ctx context.Context, key string) (*domain.Task, error)
	FindPriority(ctx context.Context) ([]domain.TaskPriority, error)
	FindStatus(ctx context.Context, opts query.QueryOptions, projectID int64) ([]domain.TaskStatus, error)
	CreateStatus(ctx context.Context, req *domain.CreateTaskStatusReq) (*domain.TaskStatus, error)
	CreateListStatus(ctx context.Context, projectID int64, req []domain.CreateListTaskStatusReq) error
	UpdateStatus(ctx context.Context, req *domain.UpdateTaskStatusReq, statusID int64) (*domain.TaskStatus, error)
	DeleteStatus(ctx context.Context, statusID int64) error
	Create(ctx context.Context, req *domain.CreateTaskReq, createdBy int64) (*domain.Task, error)
	Update(ctx context.Context, req *domain.UpdateTaskReq, taskID int64) (*domain.Task, error)
	Delete(ctx context.Context, taskID int64) error
	ReorderStatus(ctx context.Context, req *domain.ReorderTaskStatusReq, projectID int64) error
	ReorderTask(ctx context.Context, req *domain.ReorderTaskReq, projectID int64) error
	FindByToday(ctx context.Context, opts query.QueryOptions, userID int64, orgID int64) ([]domain.TaskToday, int64, error)
	FindOverdue(ctx context.Context, opts query.QueryOptions, userID int64, orgID int64) ([]domain.TaskToday, int64, error)
	CreateAttachments(ctx context.Context, req *domain.TaskAttachment, taskID int64) (*domain.TaskAttachmentFileUploadRes, error)
	FindAttachments(ctx context.Context, taskID int64) ([]domain.TaskAttachmentItem, error)
	FindAttachment(ctx context.Context, attachmentID int64) (*domain.TaskAttachmentItem, error)
	DeleteAttachment(ctx context.Context, attachmentID int64) error

	FindSubTasks(ctx context.Context, taskID int64) ([]domain.SubTask, error)
	CreateSubTask(ctx context.Context, req *domain.SubTaskReq, createdBy int64) (*domain.SubTask, error)
	UpdateSubTask(ctx context.Context, req *domain.UpdateSubTaskReq, subtaskID int64) (*domain.SubTask, error)
	DeleteSubTask(ctx context.Context, subtaskID int64) error
	ReorderSubTask(ctx context.Context, req *domain.ReorderSubTaskReq, taskID int64) error
}
