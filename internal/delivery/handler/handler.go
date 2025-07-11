package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"email-verification/internal/domain"
	"email-verification/internal/service"
)

type VerificationHandler struct {
	verificationSvc *service.VerificationService
	timeout         time.Duration
}

func NewVerificationHandler(
	verificationSvc *service.VerificationService,
	timeout time.Duration,
) *VerificationHandler {
	return &VerificationHandler{
		verificationSvc: verificationSvc,
		timeout:         timeout,
	}
}

func (h *VerificationHandler) SendVerificationCode(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	var req domain.VerificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.verificationSvc.SendVerificationCode(ctx, req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *VerificationHandler) VerifyCode(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	var req domain.VerificationCode
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.verificationSvc.VerifyCode(ctx, req.Email, req.Code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *VerificationHandler) SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/send-verification-code", h.SendVerificationCode)
	mux.HandleFunc("/api/verify-code", h.VerifyCode)
}
