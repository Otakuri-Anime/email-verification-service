package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type ElasticEmailSender struct {
	apiKey    string
	fromEmail string
	endpoint  string
}

func NewElasticEmailSender(apiKey, fromEmail, endpoint string) (*ElasticEmailSender, error) {
	return &ElasticEmailSender{
		apiKey:    apiKey,
		fromEmail: fromEmail,
		endpoint:  endpoint,
	}, nil
}

func (s *ElasticEmailSender) SendVerificationEmail(
	ctx context.Context,
	toEmail, code string,
) error {
	subject := "Your Verification Code"
	body := fmt.Sprintf("Your verification code is: %s", code)

	requestBody := map[string]string{
		"apikey":          s.apiKey,
		"from":            s.fromEmail,
		"to":              toEmail,
		"subject":         subject,
		"bodyText":        body,
		"isTransactional": "true",
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal email request: %w", err)
	}

	resp, err := http.Post(s.endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send email request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("email service returned status: %s", resp.Status)
	}

	return nil
}
