package entity

import "time"

type PageRestriction struct {
	ID          string    `json:"id"`
	PageID      string    `json:"page_id"`
	Type        string    `json:"type"`         // view | edit
	SubjectType string    `json:"subject_type"` // user | role
	SubjectID   string    `json:"subject_id"`   // user UUID or role name
	CreatedAt   time.Time `json:"created_at"`

	User *UserShort `json:"user,omitempty"`
}

type SetPageRestrictionsReq struct {
	Restrictions []PageRestrictionItem `json:"restrictions" binding:"required"`
}

type PageRestrictionItem struct {
	Type        string `json:"type"         binding:"required,oneof=view edit"`
	SubjectType string `json:"subject_type" binding:"required,oneof=user role"`
	SubjectID   string `json:"subject_id"   binding:"required"`
}

type PageAccessInfo struct {
	CanView bool `json:"can_view"`
	CanEdit bool `json:"can_edit"`
}
