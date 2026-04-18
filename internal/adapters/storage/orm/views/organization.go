package views

import (
	"gorm.io/gorm"
	"time"
	"aprilpollo/internal/core/domain"
)

type OrganizationMemberView struct {
	ID             int64
	OrganizationID int64
	UserID         int64
	Email          string
	FirstName      string
	LastName       string
	DisplayName    string
	Avatar         *string
	RoleID         int64
	StatusID       int64
	InvitedAt      *time.Time
	JoinedAt       *time.Time
	InvitedBy      *int64
	IsOwner        bool
	IsPrimary      bool

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (OrganizationMemberView) TableName() string {
	return "vw_organizations_members"
}

func (m *OrganizationMemberView) ToDomain() *domain.OrganizationMember {
	return &domain.OrganizationMember{
		ID:             m.ID,
		OrganizationID: m.OrganizationID,
		UserID:         m.UserID,
		Email:          m.Email,
		FirstName:      m.FirstName,
		LastName:       m.LastName,
		DisplayName:    m.DisplayName,
		Avatar:         m.Avatar,
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
