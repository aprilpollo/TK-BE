package domain

import (
	"time"

	"github.com/google/uuid"
)

type ProjectStatus struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type Project struct {
	ID             int64         `json:"id"`
	OrganizationID int64         `json:"organization_id"`
	Name           string        `json:"name"`
	Key            uuid.UUID     `json:"key"` // e.g., "PROJ", "WEB" for ticket prefixes
	Description    string        `json:"description,omitempty"`
	LogoURL        *string       `json:"logo_url,omitempty"`
	DueDate        *time.Time    `json:"due_date,omitempty"`
	Status         ProjectStatus `json:"status"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

type CreateProjectReq struct {
	Name           string     `json:"name" validate:"required,min=3,max=255"`
	Description    string     `json:"description,omitempty"`
	LogoURL        *string    `json:"logo_url,omitempty"`
	DueDate        *time.Time `json:"due_date,omitempty"`
}

type UpdateProjectReq struct {
	OrganizationID int64      `json:"organization_id"`
	Name           string     `json:"name" validate:"required,min=3,max=255"`
	Description    string     `json:"description,omitempty"`
	LogoURL        *string    `json:"logo_url,omitempty"`
	DueDate        *time.Time `json:"due_date,omitempty"`
	StatusID       int64      `json:"status_id,omitempty"`
}
