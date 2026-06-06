package entity

import "time"

type Workflow struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	IsDefault   bool       `json:"is_default"`
	CreatedBy   string     `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"-"`

	Statuses    []WorkflowStatus     `json:"statuses,omitempty"`
	Transitions []WorkflowTransition `json:"transitions,omitempty"`
}

type CreateWorkflowReq struct {
	Name        string  `json:"name"        validate:"required,min=2,max=255"`
	Description *string `json:"description" validate:"omitempty,max=1000"`
	IsDefault   bool    `json:"is_default"`
}

type UpdateWorkflowReq struct {
	Name        *string `json:"name"        validate:"omitempty,min=2,max=255"`
	Description *string `json:"description" validate:"omitempty,max=1000"`
	IsDefault   *bool   `json:"is_default"`
}
