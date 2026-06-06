package entity

import "time"

type FieldConfiguration struct {
	ID        string             `json:"id"`
	Name      string             `json:"name"`
	ProjectID *string            `json:"project_id,omitempty"`
	Items     []*FieldConfigItem `json:"items,omitempty"`
	CreatedAt time.Time          `json:"created_at"`
}

type FieldConfigItem struct {
	ID          string  `json:"id"`
	ConfigID    string  `json:"config_id"`
	FieldName   string  `json:"field_name"`
	IsRequired  bool    `json:"is_required"`
	IsHidden    bool    `json:"is_hidden"`
	Description *string `json:"description,omitempty"`
}

type CreateFieldConfigurationReq struct {
	Name      string                  `json:"name"       validate:"required,min=1,max=255"`
	ProjectID *string                 `json:"project_id" validate:"omitempty,uuid4"`
	Items     []CreateFieldConfigItem `json:"items"`
}

type CreateFieldConfigItem struct {
	FieldName   string  `json:"field_name"  validate:"required,max=100"`
	IsRequired  bool    `json:"is_required"`
	IsHidden    bool    `json:"is_hidden"`
	Description *string `json:"description"`
}
