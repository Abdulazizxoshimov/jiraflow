package entity

import "time"

type IssueVote struct {
	IssueID   string     `json:"issue_id"`
	UserID    string     `json:"user_id"`
	CreatedAt time.Time  `json:"created_at"`
	User      *UserShort `json:"user,omitempty"`
}

type IssueVoteSummary struct {
	Count   int          `json:"count"`
	HasMine bool         `json:"has_mine"`
	Voters  []*UserShort `json:"voters,omitempty"`
}
