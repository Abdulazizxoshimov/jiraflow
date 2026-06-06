package entity

import "time"

type IssueWatcher struct {
	IssueID   string    `json:"issue_id"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`

	User *UserShort `json:"user,omitempty"`
}
