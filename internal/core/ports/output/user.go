package output

import (
	"context"

	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/pkg/query"
)

type UserRepository interface {
	FindAll(ctx context.Context, opts query.QueryOptions) ([]domain.User, int64, error)
	FindByID(ctx context.Context, id int64) (*domain.User, error)
	Update(ctx context.Context, id int64, req *domain.UpdateUserReq) error
}
