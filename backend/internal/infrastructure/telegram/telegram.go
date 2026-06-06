package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	botToken   string
	httpClient *http.Client
}

func New(botToken string) *Client {
	return &Client{
		botToken:   botToken,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) apiURL(method string) string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/%s", c.botToken, method)
}

func (c *Client) SendMessage(ctx context.Context, chatID int64, text string) error {
	body, _ := json.Marshal(map[string]any{
		"chat_id":    chatID,
		"text":       text,
		"parse_mode": "HTML",
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.apiURL("sendMessage"), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("telegram.SendMessage build req: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("telegram.SendMessage do: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("telegram.SendMessage status %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) SetWebhook(ctx context.Context, url, secret string) error {
	body, _ := json.Marshal(map[string]any{
		"url":          url,
		"secret_token": secret,
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.apiURL("setWebhook"), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("telegram.SetWebhook build req: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("telegram.SetWebhook do: %w", err)
	}
	defer resp.Body.Close()
	return nil
}

func (c *Client) DeleteWebhook(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.apiURL("deleteWebhook"), nil)
	if err != nil {
		return fmt.Errorf("telegram.DeleteWebhook: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("telegram.DeleteWebhook do: %w", err)
	}
	defer resp.Body.Close()
	return nil
}
