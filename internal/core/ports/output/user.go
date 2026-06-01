package output

import (
	"context"

	"aprilpollo/internal/core/domain"
)

type UserRepository interface {
	FindByID(ctx context.Context, id int64) (*domain.User, error)
	Update(ctx context.Context, id int64, req *domain.UpdateUserReq) error
}
