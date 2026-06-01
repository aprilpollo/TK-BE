package input

import (
	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/pkg/query"
	"context"
)

type TaskService interface {
	List(ctx context.Context, opts query.QueryOptions, projectID int64, statusID int64) ([]domain.Task, int64, error)
	GetByKey(ctx context.Context, key string) (*domain.Task, error)
	ListPriority(ctx context.Context) ([]domain.TaskPriority, error)
	ListStatus(ctx context.Context, opts query.QueryOptions, projectID int64) ([]domain.TaskStatus, error)
	CreateStatus(ctx context.Context, req *domain.CreateTaskStatusReq) (*domain.TaskStatus, error)
	CreateListStatus(ctx context.Context, projectID int64, req []domain.CreateListTaskStatusReq) error
	UpdateStatus(ctx context.Context, req *domain.UpdateTaskStatusReq, statusID int64) (*domain.TaskStatus, error)
	DeleteStatus(ctx context.Context, statusID int64) error
	Create(ctx context.Context, req *domain.CreateTaskReq, createdBy int64) (*domain.Task, error)
	Update(ctx context.Context, req *domain.UpdateTaskReq, taskID int64) (*domain.Task, error)
	Delete(ctx context.Context, taskID int64) error
	ReorderStatus(ctx context.Context, req *domain.ReorderTaskStatusReq, projectID int64) error
	ReorderTask(ctx context.Context, req *domain.ReorderTaskReq, projectID int64) error
	ListByToday(ctx context.Context, opts query.QueryOptions, userID int64, orgID int64) ([]domain.TaskToday, int64, error)
	ListOverdue(ctx context.Context, opts query.QueryOptions, userID int64, orgID int64) ([]domain.TaskToday, int64, error)
	CreateAttachments(ctx context.Context, items []domain.UploadFileItem, taskID int64, uploadedBy int64) ([]*domain.TaskAttachmentFileUploadRes, error)
	GetAttachments(ctx context.Context, taskID int64) ([]domain.TaskAttachmentItem, error)
	DeleteAttachment(ctx context.Context, attachmentID int64) error

	ListSubTasks(ctx context.Context, taskID int64) ([]domain.SubTask, error)
	CreateSubTask(ctx context.Context, req *domain.SubTaskReq, createdBy int64) (*domain.SubTask, error)
	UpdateSubTask(ctx context.Context, req *domain.UpdateSubTaskReq, subtaskID int64) (*domain.SubTask, error)
	DeleteSubTask(ctx context.Context, subtaskID int64) error
	ReorderSubTask(ctx context.Context, req *domain.ReorderSubTaskReq, taskID int64) error
}
