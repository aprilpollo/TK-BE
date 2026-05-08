package models

import (
	"time"

	"aprilpollo/internal/core/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectModel struct {
	ID             int64     `gorm:"primaryKey"`
	OrganizationID int64     `gorm:"not null;index"`
	Name           string    `gorm:"not null;size:255"`
	Key            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Description    string    `gorm:"type:text"`
	LogoURL        *string   `gorm:"type:text"`
	StartDate      *int64
	EndDate        *int64
	StatusID       int64 `gorm:"not null;index;default:1"`

	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Relationships
	Organization *OrganizationMemberModel `gorm:"foreignKey:OrganizationID"`
	Status       *ProjectStatusModel      `gorm:"foreignKey:StatusID"`
}

type ProjectStatusModel struct {
	ID          int64  `gorm:"primaryKey;autoIncrement:false"` // manually set IDs
	Name        string `gorm:"not null;uniqueIndex;size:50"`   // active, inactive, completed, cancelled
	Description string `gorm:"type:text;size:255"`
}

type ProjectNotificationSettingModel struct {
	ProjectID int64 `gorm:"not null;uniqueIndex:uq_project_user_notif"`

	TaskAssignedEmail       bool `gorm:"default:true"`
	TaskAssignedInapp       bool `gorm:"default:true"`
	TaskStatusChangedEmail  bool `gorm:"default:false"`
	TaskStatusChangedInapp  bool `gorm:"default:true"`
	MentionedInCommentEmail bool `gorm:"default:true"`
	MentionedInCommentInapp bool `gorm:"default:true"`
	DueDateApproachingEmail bool `gorm:"default:true"`
	DueDateApproachingInapp bool `gorm:"default:true"`
	ProjectUpdatesEmail     bool `gorm:"default:false"`
	ProjectUpdatesInapp     bool `gorm:"default:true"`
	NewMemberJoinedEmail    bool `gorm:"default:false"`
	NewMemberJoinedInapp    bool `gorm:"default:false"`

	DailyDigest  bool `gorm:"default:false"`
	WeeklyDigest bool `gorm:"default:true"`

	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Project *ProjectModel `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
}

func (ProjectModel) TableName() string {
	return "projects"
}

func (ProjectStatusModel) TableName() string {
	return "project_statuses"
}

func (ProjectNotificationSettingModel) TableName() string {
	return "project_notification_settings"
}

func (m *ProjectModel) ToDomain() *domain.Project {
	status := domain.ProjectStatus{ID: m.StatusID}
	if m.Status != nil {
		status = *m.Status.ToDomain()
	}

	return &domain.Project{
		ID:             m.ID,
		OrganizationID: m.OrganizationID,
		Name:           m.Name,
		Key:            m.Key,
		Description:    m.Description,
		LogoURL:        m.LogoURL,
		StartDate:      m.StartDate,
		EndDate:        m.EndDate,
		Status:         status,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

func (m *ProjectStatusModel) ToDomain() *domain.ProjectStatus {
	return &domain.ProjectStatus{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
	}
}

func (m *ProjectNotificationSettingModel) ToDomain() *domain.ProjectNotificationSettings {
	return &domain.ProjectNotificationSettings{
		ProjectID: m.ProjectID,

		TaskAssignedEmail:       m.TaskAssignedEmail,
		TaskAssignedInapp:       m.TaskAssignedInapp,
		TaskStatusChangedEmail:  m.TaskStatusChangedEmail,
		TaskStatusChangedInapp:  m.TaskStatusChangedInapp,
		MentionedInCommentEmail: m.MentionedInCommentEmail,
		MentionedInCommentInapp: m.MentionedInCommentInapp,
		DueDateApproachingEmail: m.DueDateApproachingEmail,
		DueDateApproachingInapp: m.DueDateApproachingInapp,
		ProjectUpdatesEmail:     m.ProjectUpdatesEmail,
		ProjectUpdatesInapp:     m.ProjectUpdatesInapp,
		NewMemberJoinedEmail:    m.NewMemberJoinedEmail,
		NewMemberJoinedInapp:    m.NewMemberJoinedInapp,

		DailyDigest:  m.DailyDigest,
		WeeklyDigest: m.WeeklyDigest,

		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
