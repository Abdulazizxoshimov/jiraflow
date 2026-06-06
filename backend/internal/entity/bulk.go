package entity

// BulkUpdateIssueReq — bir vaqtda ko'p issue'ni yangilash.
type BulkUpdateIssueReq struct {
	IssueIDs    []string `json:"issue_ids"    validate:"required,min=1,dive,uuid4"`
	AssigneeID  *string  `json:"assignee_id"  validate:"omitempty,uuid4"`
	StatusID    *string  `json:"status_id"    validate:"omitempty,uuid4"`
	Priority    *string  `json:"priority"     validate:"omitempty,oneof=lowest low medium high highest"`
	SprintID    *string  `json:"sprint_id"    validate:"omitempty,uuid4"`
	LabelIDs    []string `json:"label_ids"`
	ComponentIDs []string `json:"component_ids"`
}

// BulkDeleteIssueReq — bir vaqtda ko'p issue'ni o'chirish.
type BulkDeleteIssueReq struct {
	IssueIDs []string `json:"issue_ids" validate:"required,min=1,dive,uuid4"`
}

// BulkMoveToSprintReq — bir vaqtda ko'p issue'ni sprintga ko'chirish.
type BulkMoveToSprintReq struct {
	IssueIDs []string `json:"issue_ids" validate:"required,min=1,dive,uuid4"`
	SprintID *string  `json:"sprint_id" validate:"omitempty,uuid4"` // nil => backlog
}

// BulkResult — bulk operatsiya natijasi.
type BulkResult struct {
	Updated []string `json:"updated"`
	Failed  []string `json:"failed,omitempty"`
	Total   int      `json:"total"`
}

// CloneIssueReq — issue nusxalash.
type CloneIssueReq struct {
	Title          *string `json:"title"`           // nil => original + " (copy)"
	SprintID       *string `json:"sprint_id"`       // nil => same sprint
	IncludeSubtasks bool   `json:"include_subtasks"`
	IncludeLinks    bool   `json:"include_links"`
}

// RankIssueReq — LexoRank orqali issue qayta tartiblash.
type RankIssueReq struct {
	Before *string `json:"before"` // issue ID before which to rank
	After  *string `json:"after"`  // issue ID after which to rank
}
