package entity

import "time"

type NotificationScheme struct {
	ID          string                   `json:"id"`
	Name        string                   `json:"name"`
	Description *string                  `json:"description,omitempty"`
	Rules       []*NotificationSchemeRule `json:"rules,omitempty"`
	CreatedAt   time.Time                `json:"created_at"`
}

type NotificationSchemeRule struct {
	ID            string    `json:"id"`
	SchemeID      string    `json:"scheme_id"`
	EventType     string    `json:"event_type"`
	RecipientType string    `json:"recipient_type"` // role | user | reporter | assignee | watchers
	RecipientID   *string   `json:"recipient_id,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

type CreateNotificationSchemeReq struct {
	Name        string                         `json:"name"        validate:"required,min=1,max=255"`
	Description *string                        `json:"description"`
	Rules       []CreateNotificationSchemeRule `json:"rules"`
}

type CreateNotificationSchemeRule struct {
	EventType     string  `json:"event_type"     validate:"required"`
	RecipientType string  `json:"recipient_type" validate:"required,oneof=role user reporter assignee watchers"`
	RecipientID   *string `json:"recipient_id"   validate:"omitempty,uuid4"`
}
