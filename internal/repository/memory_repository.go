package repository

import (
	"context"
	"sync"
	"time"
)

type MemoryVerificationRepository struct {
	mu    sync.RWMutex
	codes map[string]codeEntry
}

type codeEntry struct {
	code      string
	expiresAt time.Time
}

func NewMemoryVerificationRepository() *MemoryVerificationRepository {
	return &MemoryVerificationRepository{
		codes: make(map[string]codeEntry),
	}
}

func (r *MemoryVerificationRepository) StoreVerificationCode(
	ctx context.Context,
	email, code string,
	expirySeconds int,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.codes[email] = codeEntry{
		code:      code,
		expiresAt: time.Now().Add(time.Duration(expirySeconds) * time.Second),
	}
	return nil
}

func (r *MemoryVerificationRepository) GetVerificationCode(
	ctx context.Context,
	email string,
) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	entry, exists := r.codes[email]
	if !exists {
		return "", nil
	}

	if time.Now().After(entry.expiresAt) {
		delete(r.codes, email)
		return "", nil
	}

	return entry.code, nil
}

func (r *MemoryVerificationRepository) DeleteVerificationCode(
	ctx context.Context,
	email string,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.codes, email)
	return nil
}
