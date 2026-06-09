package uazapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type WhatsAppProvider struct {
	baseURL string
	token   string
	client  *http.Client
}

func NewWhatsAppProvider(baseURL, token string, client *http.Client) *WhatsAppProvider {
	if client == nil {
		client = &http.Client{Timeout: 10 * 1e9} // 10s
	}
	return &WhatsAppProvider{baseURL: strings.TrimRight(baseURL, "/"), token: token, client: client}
}

func (p *WhatsAppProvider) Send(ctx context.Context, recipient, message string) error {
	body, err := json.Marshal(map[string]string{
		"number": normalizeNumber(recipient),
		"text":   message,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/send/text", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("token", p.token)

	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("uazapi: delivery failed status=%d body=%s", resp.StatusCode, body)
		return fmt.Errorf("uazapi: unexpected status %d", resp.StatusCode)
	}
	return nil
}

// normalizeNumber strips leading '+' — uazapi expects digits only (e.g. 5511999999999).
func normalizeNumber(n string) string {
	return strings.TrimPrefix(n, "+")
}
