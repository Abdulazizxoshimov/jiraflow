package entity

import "time"

type CommentMention struct {
	CommentID string    `json:"comment_id"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`

	User *UserShort `json:"user,omitempty"`
}
