package domain

import (
	"github.com/google/uuid"
	"io"
	"time"
)

type ProjectStatus struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Project struct {
	ID             int64         `json:"id"`
	OrganizationID int64         `json:"organization_id"`
	Name           string        `json:"name"`
	Key            uuid.UUID     `json:"key"` // e.g., "PROJ", "WEB" for ticket prefixes
	Description    string        `json:"description"`
	LogoURL        *string       `json:"logo_url"`
	StartDate      *int64        `json:"start_date"`
	EndDate        *int64        `json:"end_date"`
	Status         ProjectStatus `json:"status"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

type ProjectNotificationSettings struct {
	ProjectID int64 `json:"project_id"`

	TaskAssignedEmail       bool `json:"task_assigned_email"`
	TaskAssignedInapp       bool `json:"task_assigned_inapp"`
	TaskStatusChangedEmail  bool `json:"task_status_changed_email"`
	TaskStatusChangedInapp  bool `json:"task_status_changed_inapp"`
	MentionedInCommentEmail bool `json:"mentionedincomment_email"`
	MentionedInCommentInapp bool `json:"mentionedincomment_inapp"`
	DueDateApproachingEmail bool `json:"due_date_approaching_email"`
	DueDateApproachingInapp bool `json:"due_date_approaching_inapp"`
	ProjectUpdatesEmail     bool `json:"project_updates_email"`
	ProjectUpdatesInapp     bool `json:"project_updates_inapp"`
	NewMemberJoinedEmail    bool `json:"new_member_joined_email"`
	NewMemberJoinedInapp    bool `json:"new_member_joined_inapp"`

	DailyDigest  bool `json:"daily_digest"`
	WeeklyDigest bool `json:"weekly_digest"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateProjectReq struct {
	Name        string  `json:"name" validate:"required,min=3,max=255"`
	Description string  `json:"description"`
	LogoURL     *string `json:"logo_url"`
	StartDate   *int64  `json:"start_date"`
	EndDate     *int64  `json:"end_date"`
}

type UpdateProjectReq struct {
	OrganizationID int64   `json:"organization_id"`
	Name           string  `json:"name" validate:"required,min=3,max=255"`
	Description    string  `json:"description"`
	LogoURL        *string `json:"logo_url"`
	StartDate      *int64  `json:"start_date"`
	EndDate        *int64  `json:"end_date"`
	StatusID       int64   `json:"status_id"`
}

type LogoUploadReq struct {
	File        io.Reader
	Size        int64
	ContentType string
	Filename    string
}
