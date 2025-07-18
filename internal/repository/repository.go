package repository

import "context"

type VerificationRepository interface {
	StoreVerificationCode(ctx context.Context, email, code string, expirySeconds int) error
	GetVerificationCode(ctx context.Context, email string) (string, error)
	DeleteVerificationCode(ctx context.Context, email string) error
}
