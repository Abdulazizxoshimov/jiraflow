package entity

import "time"

type AuditLog struct {
	ID         string         `json:"id"`
	UserID     *string        `json:"user_id,omitempty"`
	Action     string         `json:"action"`
	EntityType *string        `json:"entity_type,omitempty"`
	EntityID   *string        `json:"entity_id,omitempty"`
	Details    map[string]any `json:"details"`
	IPAddress  *string        `json:"ip_address,omitempty"`
	UserAgent  *string        `json:"user_agent,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`

	User *UserShort `json:"user,omitempty"`
}

type AuditLogFilter struct {
	Filter
	UserID     string `form:"user_id"     json:"user_id,omitempty"`
	Action     string `form:"action"      json:"action,omitempty"`
	EntityType string `form:"entity_type" json:"entity_type,omitempty"`
	EntityID   string `form:"entity_id"   json:"entity_id,omitempty"`
}
