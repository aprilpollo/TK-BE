package input

import (
	"context"

	"aprilpollo/internal/core/domain"
	// "aprilpollo/internal/pkg/query"

	// "github.com/google/uuid"
)

type TaskService interface {
	CreateStatus(ctx context.Context, req *domain.CreateTaskStatusReq) (*domain.TaskStatus, error)
}