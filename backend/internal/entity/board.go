package entity

import "time"

type Board struct {
	ID           string         `json:"id"`
	ProjectID    string         `json:"project_id"`
	Name         string         `json:"name"`
	Type         string         `json:"type"`          // kanban | scrum
	SwimlaneType string         `json:"swimlane_type"` // none | assignee | epic | priority | label
	Filter       map[string]any `json:"filter"`
	CreatedBy    string         `json:"created_by"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    *time.Time     `json:"-"`

	Columns []BoardColumn `json:"columns,omitempty"`
}

type CreateBoardReq struct {
	Name   string         `json:"name"   validate:"required,min=2,max=255"`
	Type   string         `json:"type"   validate:"required,oneof=kanban scrum"`
	Filter map[string]any `json:"filter"`
}

type UpdateBoardReq struct {
	Name   *string        `json:"name"   validate:"omitempty,min=2,max=255"`
	Filter map[string]any `json:"filter"`
}

type SetSwimlaneTypeReq struct {
	SwimlaneType string `json:"swimlane_type" validate:"required,oneof=none assignee epic priority label"`
}

type BoardSwimlane struct {
	Key    string   `json:"key"`    // user_id | epic_id | priority value | label_id | "none"
	Label  string   `json:"label"`  // display name
	Issues []*Issue `json:"issues"` // flat list — frontend groups by column
}

type GetBoardSwimlanesResp struct {
	SwimlaneType string           `json:"swimlane_type"`
	Columns      []*BoardColumn   `json:"columns"`
	Swimlanes    []*BoardSwimlane `json:"swimlanes"`
}
