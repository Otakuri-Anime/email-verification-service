package email

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
)

//go:embed template.html
var templateFS embed.FS

type ElasticEmailSender struct {
	apiKey        string
	fromEmail     string
	endpoint      string
	emailTmpl     *template.Template
	expiryMinutes int
}

func NewElasticEmailSender(apiKey, fromEmail, endpoint string, expiryMinutes int) (*ElasticEmailSender, error) {
	// Загружаем шаблон
	tmplContent, err := templateFS.ReadFile("template.html")
	if err != nil {
		return nil, fmt.Errorf("failed to read template: %w", err)
	}

	tmpl, err := template.New("email").Parse(string(tmplContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	return &ElasticEmailSender{
		apiKey:        apiKey,
		fromEmail:     fromEmail,
		endpoint:      endpoint,
		emailTmpl:     tmpl,
		expiryMinutes: expiryMinutes,
	}, nil
}

func (s *ElasticEmailSender) SendVerificationEmail(ctx context.Context, toEmail, code string) error {
	// Генерируем HTML из шаблона
	var body bytes.Buffer
	err := s.emailTmpl.Execute(&body, TemplateData{
		Code:          code,
		ExpiryMinutes: s.expiryMinutes,
	})
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	requestBody := map[string]interface{}{
		"Recipients": map[string][]string{
			"To": {toEmail},
		},
		"Content": map[string]interface{}{
			"From":    s.fromEmail,
			"Subject": "Твой код подтверждения",
			"Body": []map[string]string{
				{
					"ContentType": "HTML",
					"Content":     body.String(),
				},
			},
		},
		"Options": map[string]interface{}{
			"Transactional": true,
		},
	}

	// Остальная часть функции остается без изменений
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("Email JSON marshal error: %v", err)
		return err
	}

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

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error: %s", string(bodyBytes))
	}

	return nil
}
