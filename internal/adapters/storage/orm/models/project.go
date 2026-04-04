package models

import (
	"time"

	"gorm.io/gorm"
)

type ProjectModel struct {
	ID             int64    `gorm:"primaryKey"`
	OrganizationID int64   `gorm:"not null;index"`
	Name           string  `gorm:"not null;size:255"`
	Key            string  `gorm:"not null;uniqueIndex;size:255"` // e.g., "PROJ", "WEB" for ticket prefixes
	Description    string  `gorm:"type:text"`
	LogoURL        *string `gorm:"type:text"`
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
