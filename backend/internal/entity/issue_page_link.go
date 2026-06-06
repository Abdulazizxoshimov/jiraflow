package entity

import "time"

type IssuePageLink struct {
	ID        string    `json:"id"`
	IssueID   string    `json:"issue_id"`
	PageID    string    `json:"page_id"`
	LinkedBy  string    `json:"linked_by"`
	CreatedAt time.Time `json:"created_at"`

	Issue        *IssueShort `json:"issue,omitempty"`
	Page         *PageShort  `json:"page,omitempty"`
	LinkedByUser *UserShort  `json:"linked_by_user,omitempty"`
}

type IssueShort struct {
	ID      string `json:"id"`
	Key     string `json:"key"`
	Title   string `json:"title"`
	Status  string `json:"status"`
	IconURL string `json:"icon_url,omitempty"`
}

type PageShort struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	SpaceID string `json:"space_id"`
}

type CreateIssuePageLinkReq struct {
	PageID string `json:"page_id" validate:"required,uuid4"`
}
