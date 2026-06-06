package entity

import "time"

type WorkflowStatus struct {
	ID         string    `json:"id"`
	WorkflowID string    `json:"workflow_id"`
	Name       string    `json:"name"`
	Category   string    `json:"category"` // todo | in_progress | done
	Color      string    `json:"color"`
	Position   int       `json:"position"`
	IsInitial  bool      `json:"is_initial"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateWorkflowStatusReq struct {
	Name      string `json:"name"      validate:"required,min=1,max=100"`
	Category  string `json:"category"  validate:"required,oneof=todo in_progress done"`
	Color     string `json:"color"     validate:"omitempty,len=7"`
	Position  int    `json:"position"`
	IsInitial bool   `json:"is_initial"`
}

type UpdateWorkflowStatusReq struct {
	Name     *string `json:"name"     validate:"omitempty,min=1,max=100"`
	Category *string `json:"category" validate:"omitempty,oneof=todo in_progress done"`
	Color    *string `json:"color"    validate:"omitempty,len=7"`
	Position *int    `json:"position"`
}
