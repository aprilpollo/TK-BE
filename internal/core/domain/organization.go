package domain

import "time"

type Organization struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	Description  string    `json:"description"`
	LogoURL      *string   `json:"logo_url"`
	ContactEmail string    `json:"contact_email"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type OrganizationMember struct {
	ID             int64      `json:"id"`
	OrganizationID int64      `json:"organization_id"`
	UserID         int64      `json:"user_id"`
	RoleID         int64      `json:"role_id"`
	StatusID       int64      `json:"status_id"`
	IsOwner        bool       `json:"is_owner"`
	IsPrimary      bool       `json:"is_primary"`
	InvitedAt      *time.Time `json:"invited_at"`
	JoinedAt       *time.Time `json:"joined_at"`
	InvitedBy      *int64     `json:"invited_by"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type UserOrganization struct {
	Organization
	MemberID  int64      `json:"member_id"`
	RoleID    int64      `json:"role_id"`
	StatusID  int64      `json:"status_id"`
	IsOwner   bool       `json:"is_owner"`
	IsPrimary bool       `json:"is_primary"`
	JoinedAt  *time.Time `json:"joined_at"`
}

type CreateOrganizationReq struct {
	Name         string  `json:"name"`
	Slug         string  `json:"slug"`
	Description  string  `json:"description"`
	LogoURL      *string `json:"logo_url"`
	ContactEmail string  `json:"contact_email"`
}

type UpdateOrganizationReq struct {
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	LogoURL      *string `json:"logo_url"`
	ContactEmail string  `json:"contact_email"`
	IsActive     *bool   `json:"is_active"`
}

type InviteMemberReq struct {
	UserID int64 `json:"user_id"`
	RoleID int64 `json:"role_id"`
}

type UpdateMemberReq struct {
	RoleID   int64 `json:"role_id"`
	StatusID int64 `json:"status_id"`
}
