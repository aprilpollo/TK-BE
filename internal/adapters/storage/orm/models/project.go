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
	DueDate        *time.Time
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

func (ProjectModel) TableName() string {
	return "projects"
}

func (ProjectStatusModel) TableName() string {
	return "project_statuses"
}

func (m *ProjectModel) ToDomain() *domain.Project {
	return &domain.Project{
		ID:             m.ID,
		OrganizationID: m.OrganizationID,
		Name:           m.Name,
		Key:            m.Key,
		Description:    m.Description,
		LogoURL:        m.LogoURL,
		DueDate:        m.DueDate,
		Status:         *m.Status.ToDomain(),
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