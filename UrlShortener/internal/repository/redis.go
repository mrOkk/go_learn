package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
	TTL    time.Duration
}

func NewRedis(client *redis.Client, ttl time.Duration) *Redis {
	return &Redis{
		client: client,
		TTL:    ttl,
	}
}

func (r *Redis) Save(ctx context.Context, code, url string) error {
	status := r.client.Set(ctx, code, url, r.TTL)
	return status.Err()
}

func (r *Redis) Get(ctx context.Context, code string) (string, bool, error) {
	cmd := r.client.Get(ctx, code)
	if cmd.Err() != nil {
		return "", false, cmd.Err()
	}
	return cmd.Val(), true, nil
}
