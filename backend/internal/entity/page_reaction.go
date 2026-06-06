package entity

import "time"

type PageReaction struct {
	ID        string    `json:"id"`
	PageID    string    `json:"page_id"`
	UserID    string    `json:"user_id"`
	Emoji     string    `json:"emoji"`
	CreatedAt time.Time `json:"created_at"`

	User *UserShort `json:"user,omitempty"`
}

// PageReactionSummary — emoji bo'yicha guruhlangan count.
type PageReactionSummary struct {
	Emoji    string `json:"emoji"`
	Count    int    `json:"count"`
	HasMine  bool   `json:"has_mine"`
}

type ToggleReactionReq struct {
	Emoji string `json:"emoji" binding:"required,min=1,max=32"`
}
