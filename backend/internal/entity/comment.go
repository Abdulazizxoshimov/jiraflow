package entity

import "time"

type CommentReaction struct {
	ID        string    `json:"id"`
	CommentID string    `json:"comment_id"`
	UserID    string    `json:"user_id"`
	Emoji     string    `json:"emoji"`
	CreatedAt time.Time `json:"created_at"`
}

type CommentReactionSummary struct {
	Emoji      string `json:"emoji"`
	Count      int    `json:"count"`
	ReactedByMe bool  `json:"reacted_by_me"`
}

type Comment struct {
	ID          string         `json:"id"`
	ParentType  string         `json:"parent_type"` // issue | page
	ParentID    string         `json:"parent_id"`
	AuthorID    string         `json:"author_id"`
	Content     map[string]any `json:"content"`      // TipTap JSON
	ContentText string         `json:"content_text"` // plain text
	ReplyToID   *string        `json:"reply_to_id,omitempty"`
	IsEdited    bool           `json:"is_edited"`
	EditedAt    *time.Time     `json:"edited_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   *time.Time     `json:"-"`

	Author    *UserShort               `json:"author,omitempty"`
	Replies   []Comment                `json:"replies,omitempty"`
	Reactions []CommentReactionSummary `json:"reactions,omitempty"`
}

type CreateCommentReq struct {
	Content     map[string]any `json:"content"      validate:"required"`
	ContentText string         `json:"content_text" validate:"required"`
	ReplyToID   *string        `json:"reply_to_id"  validate:"omitempty,uuid4"`
}

type UpdateCommentReq struct {
	Content     map[string]any `json:"content"      validate:"required"`
	ContentText string         `json:"content_text" validate:"required"`
}
