package output

import (
	"context"

	"aprilpollo/internal/core/domain"
)

type CommentPublisher interface {
	PublishComment(ctx context.Context, taskID int64, event domain.CommentEvent) error
}
