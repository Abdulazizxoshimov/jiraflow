package entity

import "time"

type SecurityScheme struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description,omitempty"`
	ProjectID   *string          `json:"project_id,omitempty"`
	Levels      []*SecurityLevel `json:"levels,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

type SecurityLevel struct {
	ID          string                  `json:"id"`
	SchemeID    string                  `json:"scheme_id"`
	Name        string                  `json:"name"`
	Description string                  `json:"description,omitempty"`
	Members     []*SecurityLevelMember  `json:"members,omitempty"`
	CreatedAt   time.Time               `json:"created_at"`
}

type SecurityLevelMember struct {
	ID      string `json:"id"`
	LevelID string `json:"level_id"`
	Type    string `json:"type"`
	Value   string `json:"value,omitempty"`
}

type CreateSecuritySchemeReq struct {
	Name        string                    `json:"name"        validate:"required"`
	Description string                    `json:"description"`
	ProjectID   *string                   `json:"project_id"`
	Levels      []CreateSecurityLevelReq  `json:"levels"`
}

type CreateSecurityLevelReq struct {
	Name        string                         `json:"name" validate:"required"`
	Description string                         `json:"description"`
	Members     []CreateSecurityLevelMemberReq `json:"members"`
}

type CreateSecurityLevelMemberReq struct {
	Type  string `json:"type"  validate:"required,oneof=role user group reporter assignee"`
	Value string `json:"value"`
}
