package entity

import "time"

type User struct {
	ID           string     `json:"id"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	FullName     string     `json:"full_name"`
	AvatarURL    *string    `json:"avatar_url,omitempty"`
	Color        string     `json:"color"`
	Role         string     `json:"role"`
	Timezone     string     `json:"timezone"`
	Language     string     `json:"language"`
	IsActive     bool       `json:"is_active"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"-"`
}

type CreateUserReq struct {
	Email    string `json:"email"     validate:"required,email"`
	Password string `json:"password"  validate:"required,min=8,max=72"`
	FullName string `json:"full_name" validate:"required,min=2,max=255"`
	Role     string `json:"role"      validate:"omitempty,oneof=admin member viewer"`
}

type UpdateUserReq struct {
	FullName  *string `json:"full_name"  validate:"omitempty,min=2,max=255"`
	AvatarURL *string `json:"avatar_url" validate:"omitempty,url"`
	Color     *string `json:"color"      validate:"omitempty,len=7"`
	Timezone  *string `json:"timezone"   validate:"omitempty,max=64"`
	Language  *string `json:"language"   validate:"omitempty,max=8"`
}

type ChangePasswordReq struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password"     validate:"required,min=8,max=72"`
}

type UserFilter struct {
	Filter
	Role     string `form:"role"      json:"role"`
	IsActive *bool  `form:"is_active" json:"is_active,omitempty"`
}

// UserShort is a lightweight user representation used in nested responses.
type UserShort struct {
	ID        string  `json:"id"`
	FullName  string  `json:"full_name"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	Color     string  `json:"color"`
	Email     string  `json:"email"`
}
