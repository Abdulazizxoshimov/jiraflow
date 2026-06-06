package entity

import "time"

type Notification struct {
	ID         string         `json:"id"`
	UserID     string         `json:"user_id"`
	Type       string         `json:"type"`
	Payload    map[string]any `json:"payload"`
	EntityType *string        `json:"entity_type,omitempty"` // issue | page | comment
	EntityID   *string        `json:"entity_id,omitempty"`
	ActorID    *string        `json:"actor_id,omitempty"`
	ReadAt     *time.Time     `json:"read_at,omitempty"`
	EmailSentAt *time.Time    `json:"email_sent_at,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`

	Actor *UserShort `json:"actor,omitempty"`
}

type NotificationFilter struct {
	Filter
	Unread *bool `form:"unread" json:"unread,omitempty"`
}

type MarkReadReq struct {
	IDs []string `json:"ids" validate:"required,min=1"`
}
