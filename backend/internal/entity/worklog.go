package entity

import "time"

// Worklog — issue'ga sarflangan vaqt yozuvi (seconds).
type Worklog struct {
	ID          string    `json:"id"`
	IssueID     string    `json:"issue_id"`
	UserID      string    `json:"user_id"`
	TimeSpent   int       `json:"time_spent"`   // sekundda
	StartedAt   time.Time `json:"started_at"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	User *UserShort `json:"user,omitempty"`
}

type CreateWorklogReq struct {
	TimeSpent   int       `json:"time_spent"   validate:"required,gt=0"`
	StartedAt   time.Time `json:"started_at"`
	Description *string   `json:"description"  validate:"omitempty,max=2000"`
}

type UpdateWorklogReq struct {
	TimeSpent   *int      `json:"time_spent"   validate:"omitempty,gt=0"`
	StartedAt   *time.Time `json:"started_at"`
	Description *string   `json:"description"  validate:"omitempty,max=2000"`
}

type WorklogFilter struct {
	Filter
	IssueID string `form:"issue_id" json:"issue_id"`
	UserID  string `form:"user_id"  json:"user_id"`
}

// TimeSpentSummary — issue bo'yicha umumiy vaqt hisobi.
type TimeSpentSummary struct {
	IssueID            string `json:"issue_id"`
	OriginalEstimate   *int   `json:"original_estimate,omitempty"`   // sekundda
	RemainingEstimate  *int   `json:"remaining_estimate,omitempty"`  // sekundda
	TimeSpentTotal     int    `json:"time_spent_total"`              // sekundda
}
