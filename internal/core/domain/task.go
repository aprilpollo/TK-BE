package domain

import (
	"github.com/google/uuid"
)

type TaskStatus struct {
	ID          int64     `json:"id"`
	UUID        uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Color       string    `json:"color,omitempty"`
	Position    int       `json:"position,omitempty"`
}

type CreateTaskStatusReq struct {
	ProjectID   int64  `json:"project_id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Color       string `json:"color,omitempty"`
}
