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

	requestBody := map[string]interface{}{
		"Recipients": map[string][]string{
			"To": {toEmail},
		},
		"Content": map[string]interface{}{
			"From":    s.fromEmail,
			"Subject": subject,
			"Body": []map[string]string{
				{
					"ContentType": "PlainText",
					"Content":     body,
				},
			},
		},
		"Options": map[string]interface{}{
			"Transactional": true,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("Email JSON marshal error: %v", err)
		return err
	}

	log.Printf("Request payload: %s", string(jsonData))

	req, err := http.NewRequest("POST", s.endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Email request creation error: %v", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-ElasticEmail-ApiKey", s.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Email send HTTP error: %v", err)
		return err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	log.Printf("ElasticEmail API response:\nStatus: %d\nBody: %s", resp.StatusCode, string(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API error: %s", string(bodyBytes))
	}

	return nil
}
