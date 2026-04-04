package models

import (
	"time"

	"aprilpollo/internal/core/domain"
	"gorm.io/gorm"
)

type OrganizationModel struct {
	ID           int64    `gorm:"primaryKey"`
	Name         string  `gorm:"not null;size:255"`
	Slug         string  `gorm:"not null;uniqueIndex;size:100"`
	Description  string  `gorm:"type:text"`
	LogoURL      *string `gorm:"type:text"`
	ContactEmail string  `gorm:"size:255"`
	IsActive     bool    `gorm:"default:true;index"`

	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type OrganizationMemberModel struct {
	ID             int64  `gorm:"primaryKey"`
	OrganizationID int64 `gorm:"not null;uniqueIndex:idx_org_user"`
	UserID         int64 `gorm:"not null;uniqueIndex:idx_org_user;index"`
	RoleID         int64 `gorm:"not null;index"`
	StatusID       int64 `gorm:"not null;index;default:1"`
	InvitedAt      *time.Time
	JoinedAt       *time.Time
	InvitedBy      *int64 `gorm:"index"`
	IsOwner        bool   `gorm:"default:false"`
	IsPrimary      bool   `gorm:"default:false"`

	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Relationships
	Organization *OrganizationModel             `gorm:"foreignKey:OrganizationID"`
	User         *UserModel                     `gorm:"foreignKey:UserID"`
	MemberStatus *OrganizationMemberStatusModel `gorm:"foreignKey:StatusID"`
	Role         *OrganizationMemberRoleModel   `gorm:"foreignKey:RoleID"`
	Inviter      *UserModel                     `gorm:"foreignKey:InvitedBy"`
}

type OrganizationMemberStatusModel struct {
	ID          int64  `gorm:"primaryKey;autoIncrement:false"` // manually set IDs
	Name        string `gorm:"not null;uniqueIndex;size:50"`   // invited, pending, active, inactive
	Description string `gorm:"type:text;size:255"`
}

type OrganizationMemberRoleModel struct {
	ID          int64  `gorm:"primaryKey;autoIncrement:false"` // manually set IDs
	Name        string `gorm:"not null;uniqueIndex;size:50"`   // user, admin, owner
	Description string `gorm:"type:text;size:255"`
	CanView     bool   `gorm:"default:true"`
	CanCreate   bool   `gorm:"default:false"`
	CanUpdate   bool   `gorm:"default:false"`
	CanDelete   bool   `gorm:"default:false"`
}

func (OrganizationModel) TableName() string {
	return "organizations"
}

func (m *OrganizationModel) ToDomain() *domain.Organization {
	return &domain.Organization{
		ID:           m.ID,
		Name:         m.Name,
		Slug:         m.Slug,
		Description:  m.Description,
		LogoURL:      m.LogoURL,
		ContactEmail: m.ContactEmail,
		IsActive:     m.IsActive,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func FromOrganizationDomain(d *domain.Organization) *OrganizationModel {
	return &OrganizationModel{
		ID:           d.ID,
		Name:         d.Name,
		Slug:         d.Slug,
		Description:  d.Description,
		LogoURL:      d.LogoURL,
		ContactEmail: d.ContactEmail,
		IsActive:     d.IsActive,
	}
}

func (OrganizationMemberModel) TableName() string {
	return "organization_members"
}

func (m *OrganizationMemberModel) ToDomain() *domain.OrganizationMember {
	return &domain.OrganizationMember{
		ID:             m.ID,
		OrganizationID: m.OrganizationID,
		UserID:         m.UserID,
		RoleID:         m.RoleID,
		StatusID:       m.StatusID,
		IsOwner:        m.IsOwner,
		IsPrimary:      m.IsPrimary,
		InvitedAt:      m.InvitedAt,
		JoinedAt:       m.JoinedAt,
		InvitedBy:      m.InvitedBy,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

func (OrganizationMemberStatusModel) TableName() string {
	return "organization_member_statuses"
}

func (OrganizationMemberRoleModel) TableName() string {
	return "organization_member_roles"
}
