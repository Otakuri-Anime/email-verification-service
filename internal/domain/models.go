package domain

type VerificationRequest struct {
	Email string `json:"email"`
}

type VerificationResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type VerificationCode struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}
