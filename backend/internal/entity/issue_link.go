package entity

import "time"

type IssueLink struct {
	ID        string    `json:"id"`
	SourceID  string    `json:"source_id"`
	TargetID  string    `json:"target_id"`
	LinkType  string    `json:"link_type"` // relates_to | blocks | blocked_by | duplicates | is_duplicated_by
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`

	Source *Issue     `json:"source,omitempty"`
	Target *Issue     `json:"target,omitempty"`
	Creator *UserShort `json:"creator,omitempty"`
}

type CreateIssueLinkReq struct {
	TargetID string `json:"target_id" validate:"required,uuid4"`
	LinkType string `json:"link_type" validate:"required,oneof=relates_to blocks blocked_by duplicates is_duplicated_by"`
}
