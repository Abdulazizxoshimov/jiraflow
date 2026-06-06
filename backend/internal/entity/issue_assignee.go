package entity

import "time"

type IssueAssignee struct {
	IssueID   string     `json:"issue_id"`
	UserID    string     `json:"user_id"`
	IsPrimary bool       `json:"is_primary"`
	CreatedAt time.Time  `json:"created_at"`

	User *UserShort `json:"user,omitempty"`
}

type SetIssueAssigneesReq struct {
	UserIDs   []string `json:"user_ids"   binding:"required"`
	PrimaryID *string  `json:"primary_id"`
}
