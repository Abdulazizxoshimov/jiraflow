package entity

import "time"

type CustomField struct {
	ID         string         `json:"id"`
	ProjectID  string         `json:"project_id"`
	Name       string         `json:"name"`
	FieldKey   string         `json:"field_key"`
	FieldType  string         `json:"field_type"`
	IsRequired bool           `json:"is_required"`
	Options    []string       `json:"options,omitempty"`
	Position   int            `json:"position"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

type CreateCustomFieldReq struct {
	Name       string   `json:"name"       validate:"required,min=1,max=100"`
	FieldKey   string   `json:"field_key"  validate:"required,min=1,max=50"`
	FieldType  string   `json:"field_type" validate:"required,oneof=text number date select multi_select user checkbox url"`
	IsRequired bool     `json:"is_required"`
	Options    []string `json:"options"    validate:"omitempty"`
	Position   int      `json:"position"`
}

type UpdateCustomFieldReq struct {
	Name       *string  `json:"name"       validate:"omitempty,min=1,max=100"`
	IsRequired *bool    `json:"is_required"`
	Options    []string `json:"options"`
	Position   *int     `json:"position"`
}
