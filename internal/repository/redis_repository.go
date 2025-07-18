// internal/repository/redis_repository.go
package repository

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisVerificationRepository struct {
	client *redis.Client
}

func NewRedisVerificationRepository(client *redis.Client) VerificationRepository {
	return &RedisVerificationRepository{
		client: client,
	}
}

func (r *RedisVerificationRepository) StoreVerificationCode(
	ctx context.Context,
	email, code string,
	expirySeconds int,
) error {
	return r.client.Set(ctx, email, code, time.Duration(expirySeconds)*time.Second).Err()
}

func (r *RedisVerificationRepository) GetVerificationCode(
	ctx context.Context,
	email string,
) (string, error) {
	return r.client.Get(ctx, email).Result()
}

func (r *RedisVerificationRepository) DeleteVerificationCode(
	ctx context.Context,
	email string,
) error {
	return r.client.Del(ctx, email).Err()
}
