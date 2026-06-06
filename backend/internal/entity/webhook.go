package entity

import "time"

type Webhook struct {
	ID        string     `json:"id"`
	ProjectID *string    `json:"project_id,omitempty"`
	SpaceID   *string    `json:"space_id,omitempty"`
	Name      string     `json:"name"`
	URL       string     `json:"url"`
	Secret    *string    `json:"secret,omitempty"`
	Events    []string   `json:"events"`
	IsActive  bool       `json:"is_active"`
	CreatedBy string     `json:"created_by"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type WebhookDelivery struct {
	ID           string         `json:"id"`
	WebhookID    string         `json:"webhook_id"`
	Event        string         `json:"event"`
	Payload      map[string]any `json:"payload"`
	StatusCode   *int           `json:"status_code,omitempty"`
	ResponseBody *string        `json:"response_body,omitempty"`
	Success      bool           `json:"success"`
	Attempt      int            `json:"attempt"`
	ErrorMsg     *string        `json:"error_msg,omitempty"`
	DeliveredAt  time.Time      `json:"delivered_at"`
}

type CreateWebhookReq struct {
	Name      string   `json:"name"   binding:"required,min=1,max=255"`
	URL       string   `json:"url"    binding:"required,url"`
	Secret    *string  `json:"secret"`
	Events    []string `json:"events" binding:"required,min=1"`
	ProjectID *string  `json:"project_id"`
	SpaceID   *string  `json:"space_id"`
}

type UpdateWebhookReq struct {
	Name     *string  `json:"name"`
	URL      *string  `json:"url"`
	Secret   *string  `json:"secret"`
	Events   []string `json:"events"`
	IsActive *bool    `json:"is_active"`
}

// Standart webhook event nomlari
const (
	EventIssueCreated   = "issue.created"
	EventIssueUpdated   = "issue.updated"
	EventIssueDeleted   = "issue.deleted"
	EventIssueTransition = "issue.transition"
	EventCommentCreated = "comment.created"
	EventSprintStarted  = "sprint.started"
	EventSprintCompleted = "sprint.completed"
	EventPageCreated    = "page.created"
	EventPageUpdated    = "page.updated"
)
