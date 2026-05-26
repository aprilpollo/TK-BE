package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"aprilpollo/internal/core/domain"

	"github.com/redis/go-redis/v9"
)

type RedisPubSub struct {
	client *redis.Client
}

func NewRedisPubSub(r *RedisClient) *RedisPubSub {
	return &RedisPubSub{client: r.GetClient()}
}

func CommentChannel(taskID int64) string {
	return fmt.Sprintf("task:comments:%d", taskID)
}

func (p *RedisPubSub) PublishComment(ctx context.Context, taskID int64, event domain.CommentEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return p.client.Publish(ctx, CommentChannel(taskID), data).Err()
}
