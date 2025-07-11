package service

import (
	"context"
	"email-verification/internal/domain"
	"email-verification/internal/repository"
	"fmt"
	"math/rand"
	"time"
)

type EmailService interface {
	SendVerificationEmail(ctx context.Context, email, code string) error
}

type VerificationService struct {
	repo       repository.VerificationRepository
	emailSvc   EmailService
	codeLen    int
	codeExpiry time.Duration
}

func NewVerificationService(
	repo repository.VerificationRepository,
	emailSvc EmailService,
	codeLen int,
	codeExpiry time.Duration,
) *VerificationService {
	return &VerificationService{
		repo:       repo,
		emailSvc:   emailSvc,
		codeLen:    codeLen,
		codeExpiry: codeExpiry,
	}
}

func (s *VerificationService) generateCode() string {
	rand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	code := make([]byte, s.codeLen)
	for i := range code {
		code[i] = digits[rand.Intn(len(digits))]
	}
	return string(code)
}

func (s *VerificationService) SendVerificationCode(
	ctx context.Context,
	email string,
) (*domain.VerificationResponse, error) {
	code := s.generateCode()

	// save code to repository
	err := s.repo.StoreVerificationCode(ctx, email, code, int(s.codeExpiry.Seconds()))
	if err != nil {
		return nil, fmt.Errorf("failed to store verification code: %w", err)
	}

	err = s.emailSvc.SendVerificationEmail(ctx, email, code)
	if err != nil {
		return nil, fmt.Errorf("failed to send verification email: %w", err)
	}

	return &domain.VerificationResponse{
		Success: true,
		Message: "Verification code sent successfully",
	}, nil
}

func (s *VerificationService) VerifyCode(
	ctx context.Context,
	email, code string,
) (*domain.VerificationResponse, error) {
	storedCode, err := s.repo.GetVerificationCode(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get verification code: %w", err)
	}

	if storedCode == "" {
		return &domain.VerificationResponse{
			Success: false,
			Message: "Verification code expired or not found",
		}, nil
	}

	if storedCode != code {
		return &domain.VerificationResponse{
			Success: false,
			Message: "Invalid verification code",
		}, nil
	}

	err = s.repo.DeleteVerificationCode(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to delete verification code: %w", err)
	}

	return &domain.VerificationResponse{
		Success: true,
		Message: "Verification successful",
	}, nil
}
