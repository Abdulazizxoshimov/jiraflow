package entity

import "time"

// AutomationRule — loyiha uchun avtomatlashtirish qoidasi.
type AutomationRule struct {
	ID            string         `json:"id"`
	ProjectID     string         `json:"project_id"`
	Name          string         `json:"name"`
	Description   *string        `json:"description,omitempty"`
	TriggerType   string         `json:"trigger_type"` // issue.created | issue.updated | issue.transition | issue.assigned | sprint.started | sprint.completed | page.created | page.updated | scheduled
	TriggerConfig map[string]any `json:"trigger_config"`
	Conditions    []AutomationCondition `json:"conditions"`
	Actions       []AutomationAction    `json:"actions"`
	IsActive      bool           `json:"is_active"`
	CreatedBy     string         `json:"created_by"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

// AutomationCondition — qoida ishlaydigan shart.
type AutomationCondition struct {
	Field    string `json:"field"`    // type | priority | status | assignee | label | sprint
	Operator string `json:"operator"` // = | != | in | not_in | is_empty | is_not_empty
	Value    any    `json:"value"`
}

// AutomationAction — shart to'g'ri bo'lganda bajariladigan harakat.
type AutomationAction struct {
	Type   string         `json:"type"`   // transition_issue | assign_issue | add_label | create_issue | send_notification | set_field
	Config map[string]any `json:"config"` // action-specific params
}

// AutomationLog — qoida ishlanganining tarixi.
type AutomationLog struct {
	ID         string     `json:"id"`
	RuleID     string     `json:"rule_id"`
	EntityID   *string    `json:"entity_id,omitempty"`
	EntityType *string    `json:"entity_type,omitempty"`
	Status     string     `json:"status"` // success | failed | skipped
	ExecutedAt time.Time  `json:"executed_at"`
	ErrorMsg   *string    `json:"error_msg,omitempty"`
}

// AutomationEvent — dispatcher tomonidan yuboriladi.
type AutomationEvent struct {
	Type      string         // trigger_type bilan mos
	ProjectID string
	EntityID  string
	EntityType string
	Payload   map[string]any // issue, sprint, page ma'lumotlari
}

// CreateAutomationRuleReq — yangi qoida yaratish so'rovi.
type CreateAutomationRuleReq struct {
	Name          string                `json:"name"         validate:"required,min=1,max=255"`
	Description   *string               `json:"description"`
	TriggerType   string                `json:"trigger_type" validate:"required,oneof=issue.created issue.updated issue.transition issue.assigned sprint.started sprint.completed page.created page.updated scheduled"`
	TriggerConfig map[string]any        `json:"trigger_config"`
	Conditions    []AutomationCondition `json:"conditions"`
	Actions       []AutomationAction    `json:"actions"      validate:"required,min=1"`
}

// UpdateAutomationRuleReq — qoidani yangilash so'rovi.
type UpdateAutomationRuleReq struct {
	Name          *string               `json:"name"         validate:"omitempty,min=1,max=255"`
	Description   *string               `json:"description"`
	TriggerType   *string               `json:"trigger_type" validate:"omitempty,oneof=issue.created issue.updated issue.transition issue.assigned sprint.started sprint.completed page.created page.updated scheduled"`
	TriggerConfig map[string]any        `json:"trigger_config"`
	Conditions    []AutomationCondition `json:"conditions"`
	Actions       []AutomationAction    `json:"actions"      validate:"omitempty,min=1"`
	IsActive      *bool                 `json:"is_active"`
}

// AutomationFilter — qoidalar ro'yxati filtri.
type AutomationFilter struct {
	Filter
	ProjectID   string `form:"project_id"`
	TriggerType string `form:"trigger_type"`
	IsActive    *bool  `form:"is_active"`
}
