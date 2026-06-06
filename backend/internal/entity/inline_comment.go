package entity

import "time"

type InlineComment struct {
	ID         string     `json:"id"`
	PageID     string     `json:"page_id"`
	AuthorID   string     `json:"author_id"`
	AnchorID   string     `json:"anchor_id"`
	QuoteText  *string    `json:"quote_text,omitempty"`
	Body       string     `json:"body"`
	Resolved   bool       `json:"resolved"`
	ResolvedBy *string    `json:"resolved_by,omitempty"`
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"-"`

	Author *UserShort `json:"author,omitempty"`
}

type CreateInlineCommentReq struct {
	AnchorID  string  `json:"anchor_id"  binding:"required"`
	QuoteText *string `json:"quote_text"`
	Body      string  `json:"body"       binding:"required,min=1"`
}

type UpdateInlineCommentReq struct {
	Body string `json:"body" binding:"required,min=1"`
}
