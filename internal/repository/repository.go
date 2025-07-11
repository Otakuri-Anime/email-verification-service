package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisVerificationRepository struct {
	client *redis.Client
}

func NewRedisVerificationRepository(client *redis.Client) *RedisVerificationRepository {
	return &RedisVerificationRepository{client: client}
}

func (r *RedisVerificationRepository) StoreVerificationCode(
	ctx context.Context,
	email, code string,
	expirySeconds int,
) error {
	key := "verification:" + email
	err := r.client.Set(ctx, key, code, time.Duration(expirySeconds)*time.Second).Err()
	if err != nil {
		return fmt.Errorf("failed to store verification code: %w", err)
	}
	return nil
}

func (r *RedisVerificationRepository) GetVerificationCode(
	ctx context.Context,
	email string,
) (string, error) {
	key := "verification:" + email
	code, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("failed to get verification code: %w", err)
	}
	return code, nil
}

func (r *RedisVerificationRepository) DeleteVerificationCode(
	ctx context.Context,
	email string,
) error {
	key := "verification:" + email
	_, err := r.client.Del(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to delete verification code: %w", err)
	}
	return nil
}
