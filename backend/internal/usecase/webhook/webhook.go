package webhook

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
)

type useCase struct {
	repo   repository.WebhookRepository
	log    logger.Logger
	client *http.Client
}

func New(repo repository.WebhookRepository, log logger.Logger) UseCase {
	return &useCase{
		repo: repo,
		log:  log,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (uc *useCase) Create(ctx context.Context, actorID string, req *entity.CreateWebhookReq) (*entity.Webhook, error) {
	wh := &entity.Webhook{
		Name:      req.Name,
		URL:       req.URL,
		Secret:    req.Secret,
		Events:    req.Events,
		IsActive:  true,
		CreatedBy: actorID,
		ProjectID: req.ProjectID,
		SpaceID:   req.SpaceID,
	}
	if err := uc.repo.Create(ctx, wh); err != nil {
		return nil, fmt.Errorf("webhook.Create: %w", err)
	}
	return wh, nil
}

func (uc *useCase) GetByID(ctx context.Context, id string) (*entity.Webhook, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *useCase) ListByProject(ctx context.Context, projectID string) ([]*entity.Webhook, error) {
	return uc.repo.ListByProject(ctx, projectID)
}

func (uc *useCase) ListBySpace(ctx context.Context, spaceID string) ([]*entity.Webhook, error) {
	return uc.repo.ListBySpace(ctx, spaceID)
}

func (uc *useCase) Update(ctx context.Context, id string, req *entity.UpdateWebhookReq) (*entity.Webhook, error) {
	wh, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		wh.Name = *req.Name
	}
	if req.URL != nil {
		wh.URL = *req.URL
	}
	if req.Secret != nil {
		wh.Secret = req.Secret
	}
	if req.Events != nil {
		wh.Events = req.Events
	}
	if req.IsActive != nil {
		wh.IsActive = *req.IsActive
	}
	if err := uc.repo.Update(ctx, wh); err != nil {
		return nil, fmt.Errorf("webhook.Update: %w", err)
	}
	return wh, nil
}

func (uc *useCase) Delete(ctx context.Context, id string) error {
	if _, err := uc.repo.GetByID(ctx, id); err != nil {
		return err
	}
	return uc.repo.Delete(ctx, id)
}

func (uc *useCase) Trigger(ctx context.Context, event string, projectID, spaceID *string, payload map[string]any) error {
	hooks, err := uc.repo.FindByEvent(ctx, event, projectID, spaceID)
	if err != nil {
		return fmt.Errorf("webhook.Trigger FindByEvent: %w", err)
	}

	for _, wh := range hooks {
		go uc.deliver(wh, event, payload)
	}
	return nil
}

const maxRetries = 3

// retryDelays — exponential backoff: 1s, 5s, 30s
var retryDelays = []time.Duration{1 * time.Second, 5 * time.Second, 30 * time.Second}

func (uc *useCase) deliver(wh *entity.Webhook, event string, payload map[string]any) {
	uc.deliverWithRetry(wh, event, payload, 1)
}

func (uc *useCase) deliverWithRetry(wh *entity.Webhook, event string, payload map[string]any, attempt int) {
	ctx := context.Background()
	d := &entity.WebhookDelivery{
		WebhookID: wh.ID,
		Event:     event,
		Payload:   payload,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		uc.log.Error(ctx, "webhook.deliver marshal", logger.SafeString("err", err.Error()))
		d.Success = false
		errMsg := err.Error()
		d.ErrorMsg = &errMsg
		_ = uc.repo.SaveDelivery(ctx, d)
		return
	}

	req, err := http.NewRequest(http.MethodPost, wh.URL, bytes.NewReader(body))
	if err != nil {
		d.Success = false
		errMsg := err.Error()
		d.ErrorMsg = &errMsg
		_ = uc.repo.SaveDelivery(ctx, d)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-JiraFlow-Event", event)
	req.Header.Set("X-JiraFlow-Attempt", fmt.Sprintf("%d", attempt))

	if wh.Secret != nil && *wh.Secret != "" {
		sig := sign(body, *wh.Secret)
		req.Header.Set("X-JiraFlow-Signature", "sha256="+sig)
	}

	resp, err := uc.client.Do(req)
	if err != nil {
		errMsg := err.Error()
		d.ResponseBody = &errMsg
		d.ErrorMsg = &errMsg
		d.Success = false
		d.Attempt = attempt
		_ = uc.repo.SaveDelivery(ctx, d)

		if attempt < maxRetries {
			delay := retryDelays[attempt-1]
			uc.log.Warn(ctx, "webhook.deliver failed, retrying",
				logger.String("webhook_id", wh.ID),
				logger.String("delay", delay.String()),
			)
			time.AfterFunc(delay, func() {
				uc.deliverWithRetry(wh, event, payload, attempt+1)
			})
		}
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	respStr := string(respBody)
	d.StatusCode = &resp.StatusCode
	d.ResponseBody = &respStr
	d.Success = resp.StatusCode >= 200 && resp.StatusCode < 300
	d.Attempt = attempt

	if !d.Success && attempt < maxRetries {
		delay := retryDelays[attempt-1]
		uc.log.Warn(ctx, "webhook.deliver non-2xx, retrying",
			logger.String("webhook_id", wh.ID),
			logger.SafeString("status", fmt.Sprintf("%d", resp.StatusCode)),
			logger.String("delay", delay.String()),
		)
		_ = uc.repo.SaveDelivery(ctx, d)
		time.AfterFunc(delay, func() {
			uc.deliverWithRetry(wh, event, payload, attempt+1)
		})
		return
	}

	_ = uc.repo.SaveDelivery(ctx, d)
}

func sign(body []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}

func (uc *useCase) ListDeliveries(ctx context.Context, webhookID string, limit int) ([]*entity.WebhookDelivery, error) {
	return uc.repo.ListDeliveries(ctx, webhookID, limit)
}
