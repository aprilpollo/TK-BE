package domain

import (
	"aprilpollo/internal/utils"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          int64     `json:"id"`
	ProjectID   int64     `json:"project_id"`
	Key         uuid.UUID `json:"key"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StatusID    int64     `json:"status_id"`
	PriorityID  int64     `json:"priority_id"`
	ParentID    *int64    `json:"parent_id"`
	Position    int       `json:"position"`
	StartDate   *int64    `json:"start_date"`
	EndDate     *int64    `json:"end_date"`
	AllDay      bool      `json:"all_day"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Status   TaskStatus   `json:"status"`
	Priority TaskPriority `json:"priority"`
	Assigns  []TaskAssign `json:"assignees"`
}

type TaskAssign struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
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

type CreateListTaskStatusReq struct {
	UUID        *uuid.UUID `json:"uuid"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Color       string     `json:"color"`
}

type ReqReorderTaskStatus struct {
	Updates []ReorderTaskStatus `json:"updates"`
}
type ReorderTaskStatus struct {
	ID       int64 `json:"id"`
	Position int   `json:"position"`
}

type ReqReorderTask struct {
	Updates []ReorderTask `json:"updates"`
}
type ReorderTask struct {
	ID       int64  `json:"id"`
	Position int    `json:"position"`
	StatusID string `json:"status_id"` // UUID of the new status
}

type TaskReq struct {
	ProjectID   int64   `json:"project_id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	StatusID    int64   `json:"status_id"`
	PriorityID  int64   `json:"priority_id"`
	ParentID    *int64  `json:"parent_id"`
	StartDate   *int64  `json:"start_date"`
	EndDate     *int64  `json:"end_date"`
	AllDay      *bool   `json:"all_day"`
	AssigneeIDs []int64 `json:"assignee_ids"`
}

type UpdateTaskStatusReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

type UpdateTaskReq struct {
	Title       string                `json:"title"`
	Description string                `json:"description"`
	StatusID    int64                 `json:"status_id"`
	PriorityID  int64                 `json:"priority_id"`
	StartDate   utils.Nullable[int64] `json:"start_date"`
	EndDate     utils.Nullable[int64] `json:"end_date"`
	AllDay      *bool                 `json:"all_day"`
}
