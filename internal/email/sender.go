package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
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
		log.Printf("Email JSON marshal error: %v", err)
		return err
	}

	// Логируем сам запрос (без API ключа)
	log.Printf("Sending email to: %s\nRequest: %+v", toEmail, map[string]string{
		"from":    s.fromEmail,
		"to":      toEmail,
		"subject": subject,
	})

	resp, err := http.Post(s.endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Email send HTTP error: %v", err)
		return err
	}
	defer resp.Body.Close()

	// Читаем полный ответ
	bodyBytes, _ := io.ReadAll(resp.Body)
	log.Printf("ElasticEmail API response:\nStatus: %d\nBody: %s", resp.StatusCode, string(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API error: %s", string(bodyBytes))
	}

	return nil
}
