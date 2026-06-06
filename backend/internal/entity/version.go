package entity

import "time"

// Version — loyiha versiyasi / release (v1.0, v2.3-beta va h.k.).
type Version struct {
	ID          string     `json:"id"`
	ProjectID   string     `json:"project_id"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	Status      string     `json:"status"` // unreleased | released | archived
	StartDate   *time.Time `json:"start_date,omitempty"`
	ReleaseDate *time.Time `json:"release_date,omitempty"`
	ReleasedAt  *time.Time `json:"released_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"-"`

	// Computed fields
	IssueCount  int `json:"issue_count,omitempty"`
	DoneCount   int `json:"done_count,omitempty"`
	Progress    int `json:"progress,omitempty"` // 0-100 %
}

type CreateVersionReq struct {
	Name        string     `json:"name"         validate:"required,min=1,max=100"`
	Description *string    `json:"description"  validate:"omitempty,max=2000"`
	StartDate   *time.Time `json:"start_date"`
	ReleaseDate *time.Time `json:"release_date"`
}

type UpdateVersionReq struct {
	Name        *string    `json:"name"         validate:"omitempty,min=1,max=100"`
	Description *string    `json:"description"  validate:"omitempty,max=2000"`
	StartDate   *time.Time `json:"start_date"`
	ReleaseDate *time.Time `json:"release_date"`
}

type ReleaseVersionReq struct {
	ReleasedAt *time.Time `json:"released_at"` // nil => NOW()
}
