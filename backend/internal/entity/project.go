package entity

import "time"

type Project struct {
	ID           string     `json:"id"`
	Key          string     `json:"key"`
	Name         string     `json:"name"`
	Description  *string    `json:"description,omitempty"`
	IconURL      *string    `json:"icon_url,omitempty"`
	LeadID       string     `json:"lead_id"`
	WorkflowID   string     `json:"workflow_id"`
	IssueCounter int64      `json:"issue_counter"`
	IsArchived   bool       `json:"is_archived"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"-"`

	Lead     *UserShort `json:"lead,omitempty"`
	Workflow *Workflow  `json:"workflow,omitempty"`
}

type CreateProjectReq struct {
	Key        string  `json:"key"         validate:"required,min=2,max=10"`
	Name       string  `json:"name"        validate:"required,min=2,max=255"`
	Description *string `json:"description" validate:"omitempty,max=2000"`
	IconURL    *string `json:"icon_url"    validate:"omitempty,url"`
	WorkflowID string  `json:"workflow_id" validate:"required,uuid4"`
}

type UpdateProjectReq struct {
	Name        *string `json:"name"        validate:"omitempty,min=2,max=255"`
	Description *string `json:"description" validate:"omitempty,max=2000"`
	IconURL     *string `json:"icon_url"    validate:"omitempty,url"`
	LeadID      *string `json:"lead_id"     validate:"omitempty,uuid4"`
	WorkflowID  *string `json:"workflow_id" validate:"omitempty,uuid4"`
	IsArchived  *bool   `json:"is_archived"`
}

type ProjectFilter struct {
	Filter
	IsArchived *bool  `form:"is_archived" json:"is_archived,omitempty"`
	LeadID     string `form:"lead_id"     json:"lead_id,omitempty"`
}

type ProjectDashboard struct {
	OpenIssues          int                    `json:"open_issues"`
	ClosedIssues        int                    `json:"closed_issues"`
	OverdueIssues       int                    `json:"overdue_issues"`
	TotalIssues         int                    `json:"total_issues"`
	ActiveSprint        *Sprint                `json:"active_sprint,omitempty"`
	IssuesByPriority    map[string]int         `json:"issues_by_priority"`
	IssuesByType        map[string]int         `json:"issues_by_type"`
	IssuesByStatus      map[string]int         `json:"issues_by_status"`
	IssuesByAssignee    []AssigneeDistribution `json:"issues_by_assignee"`
	RecentActivityCount int                    `json:"recent_activity_count"`
}

type AssigneeDistribution struct {
	UserID   string `json:"user_id"`
	FullName string `json:"full_name"`
	Count    int    `json:"count"`
}
