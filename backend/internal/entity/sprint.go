package entity

import "time"

type Sprint struct {
	ID          string     `json:"id"`
	ProjectID   string     `json:"project_id"`
	Name        string     `json:"name"`
	Goal        *string    `json:"goal,omitempty"`
	Status      string     `json:"status"` // planned | active | completed
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	CreatedBy   string     `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"-"`
}

type CreateSprintReq struct {
	Name      string     `json:"name"       validate:"required,min=2,max=255"`
	Goal      *string    `json:"goal"       validate:"omitempty,max=2000"`
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
}

type UpdateSprintReq struct {
	Name      *string    `json:"name"       validate:"omitempty,min=2,max=255"`
	Goal      *string    `json:"goal"       validate:"omitempty,max=2000"`
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
}

type SprintFilter struct {
	Filter
	Status string `form:"status" json:"status" validate:"omitempty,oneof=planned active completed"`
}

// SprintPlanningView is returned by GET /projects/:id/sprint-planning.
type SprintPlanningView struct {
	ActiveSprint *Sprint  `json:"active_sprint,omitempty"`
	BacklogItems []*Issue `json:"backlog_items"`
	SprintItems  []*Issue `json:"sprint_items"`
}

// AssignToSprintReq is the body for POST /projects/:id/sprint-planning.
type AssignToSprintReq struct {
	SprintID string   `json:"sprint_id" validate:"required,uuid4"`
	IssueIDs []string `json:"issue_ids" validate:"required,min=1"`
}

// SprintCapacity holds story-point totals per assignee.
type SprintCapacity struct {
	TotalPoints  int                      `json:"total_points"`
	ByAssignee   []AssigneeCapacity       `json:"by_assignee"`
}

type AssigneeCapacity struct {
	UserID      string `json:"user_id"`
	DisplayName string `json:"display_name"`
	Points      int    `json:"points"`
}

// UpdateSprintGoalReq updates just the sprint goal.
type UpdateSprintGoalReq struct {
	Goal string `json:"goal" validate:"required,max=2000"`
}
