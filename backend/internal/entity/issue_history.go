package entity

import "time"

type IssueHistory struct {
	ID        string         `json:"id"`
	IssueID   string         `json:"issue_id"`
	UserID    *string        `json:"user_id,omitempty"`
	Field     string         `json:"field"`
	OldValue  map[string]any `json:"old_value,omitempty"`
	NewValue  map[string]any `json:"new_value,omitempty"`
	CreatedAt time.Time      `json:"created_at"`

	User *UserShort `json:"user,omitempty"`
}
