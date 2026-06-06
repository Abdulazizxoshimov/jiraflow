package entity

import "time"

type WorkflowTransition struct {
	ID           string    `json:"id"`
	WorkflowID   string    `json:"workflow_id"`
	FromStatusID *string   `json:"from_status_id"` // nil = from any status
	ToStatusID   string    `json:"to_status_id"`
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"created_at"`

	FromStatus *WorkflowStatus `json:"from_status,omitempty"`
	ToStatus   *WorkflowStatus `json:"to_status,omitempty"`
}

type CreateWorkflowTransitionReq struct {
	FromStatusID *string `json:"from_status_id"`
	ToStatusID   string  `json:"to_status_id" validate:"required,uuid4"`
	Name         string  `json:"name"         validate:"required,min=1,max=100"`
}
