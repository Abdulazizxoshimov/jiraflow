package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type Client struct {
	botToken   string
	httpClient *http.Client
}

func New(botToken string) *Client {
	return &Client{
		botToken: botToken,
		// No global timeout — each request sets its own context deadline.
		// A global timeout would race with the 30s long-poll context and
		// cause spurious cancellations.
		httpClient: &http.Client{},
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
	body, _ := json.Marshal(map[string]any{"drop_pending_updates": false})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.apiURL("deleteWebhook"), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("telegram.DeleteWebhook: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("telegram.DeleteWebhook do: %w", err)
	}
	defer resp.Body.Close()
	return nil
}

type getUpdatesResp struct {
	OK     bool                    `json:"ok"`
	Result []entity.TelegramUpdate `json:"result"`
}

// Poll long-polls Telegram for updates and calls handler for each one.
// Blocks until ctx is cancelled. Suitable for local development without a public URL.
func (c *Client) Poll(ctx context.Context, handler func(context.Context, *entity.TelegramUpdate)) error {
	// First delete any existing webhook so polling works
	_ = c.DeleteWebhook(ctx)

	offset := int64(0)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		updates, newOffset, err := c.getUpdates(ctx, offset, 30)
		if err != nil {
			// On network error, wait briefly and retry
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(3 * time.Second):
			}
			continue
		}
		offset = newOffset
		for i := range updates {
			handler(ctx, &updates[i])
		}
	}
}

func (c *Client) getUpdates(ctx context.Context, offset, timeout int64) ([]entity.TelegramUpdate, int64, error) {
	body, _ := json.Marshal(map[string]any{
		"offset":  offset,
		"timeout": timeout,
		"limit":   100,
	})
	// Use a longer HTTP timeout than the long-poll timeout
	pollCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout+5)*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(pollCtx, http.MethodPost, c.apiURL("getUpdates"), bytes.NewReader(body))
	if err != nil {
		return nil, offset, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, offset, err
	}
	defer resp.Body.Close()

	var result getUpdatesResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, offset, err
	}
	if !result.OK || len(result.Result) == 0 {
		return nil, offset, nil
	}

	// Next offset = last update_id + 1
	lastID := result.Result[len(result.Result)-1].UpdateID
	return result.Result, lastID + 1, nil
}
