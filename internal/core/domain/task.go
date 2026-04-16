package domain

import (
	"github.com/google/uuid"
	"time"
)

type Task struct {
	ID          int64      `json:"id"`
	ProjectID   int64      `json:"project_id"`
	Key         uuid.UUID  `json:"key"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	StatusID    int64      `json:"status_id"`
	PriorityID  int64      `json:"priority_id"`
	ParentID    *int64     `json:"parent_id"`
	Position    int        `json:"position"`
	DueDate     *time.Time `json:"due_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	Status      TaskStatus   `json:"status"`
	Priority    TaskPriority `json:"priority"`
}

type TaskStatus struct {
	ID          int64     `json:"id"`
	UUID        uuid.UUID `json:"uuid"`
	ProjectID   int64     `json:"project_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Color       string    `json:"color"`
	Position    int       `json:"position"`
}

type TaskPriority struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

type CreateTaskStatusReq struct {
	ProjectID   int64  `json:"project_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

type ReqReorderTaskStatus struct {
	Updates []ReqReorderItem `json:"updates"`
}
type ReqReorderItem struct {
	ID       int64 `json:"id"`
	Position int   `json:"position"`
}
