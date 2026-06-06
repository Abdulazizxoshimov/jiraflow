package entity

import "time"

// SprintReport — sprint yakuniy hisoboti.
type SprintReport struct {
	SprintID       string    `json:"sprint_id"`
	SprintName     string    `json:"sprint_name"`
	StartDate      *time.Time `json:"start_date,omitempty"`
	EndDate        *time.Time `json:"end_date,omitempty"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`

	TotalIssues     int `json:"total_issues"`
	CompletedIssues int `json:"completed_issues"`
	IncompleteIssues int `json:"incomplete_issues"`
	AddedDuringSprint int `json:"added_during_sprint"`

	TotalStoryPoints     int `json:"total_story_points"`
	CompletedStoryPoints int `json:"completed_story_points"`

	CompletionRate float64 `json:"completion_rate"` // 0.0–1.0

	CompletedList  []*Issue `json:"completed_list,omitempty"`
	IncompleteList []*Issue `json:"incomplete_list,omitempty"`
}

// BurndownPoint — burndown chart uchun bir kun qiymati.
type BurndownPoint struct {
	Date            time.Time `json:"date"`
	RemainingPoints int       `json:"remaining_points"`
	IdealPoints     float64   `json:"ideal_points"`
}

// BurndownChart — sprint burndown grafigi.
type BurndownChart struct {
	SprintID    string          `json:"sprint_id"`
	TotalPoints int             `json:"total_points"`
	Points      []BurndownPoint `json:"points"`
}

// VelocityPoint — bitta sprint uchun velocity ma'lumoti.
type VelocityPoint struct {
	SprintID    string  `json:"sprint_id"`
	SprintName  string  `json:"sprint_name"`
	Committed   int     `json:"committed"`   // sprint boshidagi story points
	Completed   int     `json:"completed"`   // bajarilgan story points
}

// VelocityReport — loyiha bo'yicha velocity tarixi.
type VelocityReport struct {
	ProjectID string          `json:"project_id"`
	Sprints   []VelocityPoint `json:"sprints"`
	Average   float64         `json:"average_velocity"`
}

// BurnupPoint — burnup chart uchun bir kun qiymati.
type BurnupPoint struct {
	Date           time.Time `json:"date"`
	CompletedPoints int      `json:"completed_points"`
	TotalScope      int      `json:"total_scope"` // sprint scope that day
}

// BurnupChart — sprint burnup grafigi (scope vs completed).
type BurnupChart struct {
	SprintID    string        `json:"sprint_id"`
	TotalPoints int           `json:"total_points"`
	Points      []BurnupPoint `json:"points"`
}

// CFDPoint — Cumulative Flow Diagram uchun bir kun qiymati.
type CFDPoint struct {
	Date     time.Time      `json:"date"`
	Counts   map[string]int `json:"counts"` // statusID → issue count
}

// CFDChart — Cumulative Flow Diagram grafigi.
type CFDChart struct {
	ProjectID string     `json:"project_id"`
	Statuses  []string   `json:"statuses"`
	Points    []CFDPoint `json:"points"`
}

// RoadmapItem — roadmap ko'rinishi uchun epic ma'lumoti.
type RoadmapItem struct {
	ID          string     `json:"id"`
	IssueNumber int        `json:"issue_number"`
	Title       string     `json:"title"`
	StatusID    string     `json:"status_id"`
	Priority    string     `json:"priority"`
	AssigneeID  *string    `json:"assignee_id,omitempty"`
	StartDate   *time.Time `json:"start_date,omitempty"` // sprint.start_date yoki due_date
	EndDate     *time.Time `json:"end_date,omitempty"`   // sprint.end_date yoki due_date
	Progress    float64    `json:"progress"`             // 0–100

	Status   *WorkflowStatus `json:"status,omitempty"`
	Assignee *UserShort      `json:"assignee,omitempty"`
	Children []*RoadmapItem  `json:"children,omitempty"` // stories/tasks under epic
}
