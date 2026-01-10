package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type BrevoMailer struct {
	APIKey string
}

func (b *BrevoMailer) Send(to, subject, body string) error {
	payload := map[string]interface{}{
		"sender": map[string]string{
			"name":  "Time Manager",
			"email": "lucasnimes30000@gmail.com",
		},
		"to": []map[string]string{
			{"email": to},
		},
		"subject":     subject,
		"htmlContent": body,
	}

	data, _ := json.Marshal(payload)

	req, err := http.NewRequest(
		"POST",
		"https://api.brevo.com/v3/smtp/email",
		bytes.NewBuffer(data),
	)
	if err != nil {
		return err
	}

	req.Header.Set("api-key", b.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("brevo mail error: %s", resp.Status)
	}

	return nil
}
