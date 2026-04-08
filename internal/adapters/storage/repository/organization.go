package repository

import (
	"context"
	"errors"
	"time"

	"aprilpollo/internal/adapters/storage/orm/models"
	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/core/ports/output"
	"aprilpollo/internal/pkg/query"
	"aprilpollo/internal/pkg/query/gormq"
	"aprilpollo/internal/utils"

	"gorm.io/gorm"
)

type organizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) output.OrganizationRepository {
	return &organizationRepository{db: db}
}

func (r *organizationRepository) FindAll(ctx context.Context, opts query.QueryOptions) ([]domain.Organization, int64, error) {
	var rows []models.OrganizationModel
	var total int64

	base := r.db.WithContext(ctx).Model(&models.OrganizationModel{})

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := gormq.ApplyToGorm(r.db.WithContext(ctx).Model(&models.OrganizationModel{}), opts).Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	orgs := make([]domain.Organization, len(rows))
	for i, row := range rows {
		orgs[i] = *row.ToDomain()
	}

	return orgs, total, nil
}

func (r *organizationRepository) FindByID(ctx context.Context, id int64) (*domain.Organization, error) {
	var row models.OrganizationModel
	if err := r.db.WithContext(ctx).First(&row, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return row.ToDomain(), nil
}

func (r *organizationRepository) CreateWithOwner(ctx context.Context, org *domain.Organization, ownerUserID int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		row := models.FromOrganizationDomain(org)
		if err := tx.Create(row).Error; err != nil {
			return err
		}
		*org = *row.ToDomain()

		now := time.Now()
		member := models.OrganizationMemberModel{
			OrganizationID: org.ID,
			UserID:         ownerUserID,
			RoleID:         1, // owner
			StatusID:       1, // active
			IsOwner:        true,
			IsPrimary:      true,
			JoinedAt:       &now,
		}
		return tx.Create(&member).Error
	})
}

func (r *organizationRepository) Update(ctx context.Context, id int64, req *domain.UpdateOrganizationReq) error {
	return r.db.WithContext(ctx).Model(&models.OrganizationModel{}).Where("id = ?", id).
		Updates(utils.StructToMap(req)).Error
}

func (r *organizationRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.OrganizationModel{}, id).Error
}

func (r *organizationRepository) FindByUserID(ctx context.Context, userID int64, opts query.QueryOptions) ([]domain.UserOrganization, int64, error) {
	var rows []models.OrganizationMemberModel
	var total int64

	base := r.db.WithContext(ctx).Model(&models.OrganizationMemberModel{}).
		Where("user_id = ?", userID)

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := gormq.ApplyToGorm(base, opts).Preload("Organization").Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	result := make([]domain.UserOrganization, 0, len(rows))
	for _, row := range rows {
		if row.Organization == nil {
			continue
		}
		result = append(result, domain.UserOrganization{
			Organization: *row.Organization.ToDomain(),
			MemberID:     row.ID,
			RoleID:       row.RoleID,
			StatusID:     row.StatusID,
			IsOwner:      row.IsOwner,
			IsPrimary:    row.IsPrimary,
			JoinedAt:     row.JoinedAt,
		})
	}

	return result, total, nil
}

func (r *organizationRepository) FindByUserIDWithPrimaryDetails(ctx context.Context, userID int64) ([]domain.UserOrganizationWithDetail, error) {
	var rows []models.OrganizationMemberModel

	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("Organization").
		Preload("Role").
		Find(&rows).Error; err != nil {
		return nil, err
	}

	// collect roleID of primary org to fetch page permissions
	var primaryRoleID int64
	for _, row := range rows {
		if row.IsPrimary {
			primaryRoleID = row.RoleID
			break
		}
	}

	var permMap map[int64][]domain.OrganizationMemberPagePermission
	if primaryRoleID > 0 {
		var perms []models.OrganizationMemberPagePermissionModel
		if err := r.db.WithContext(ctx).
			Where("role_id = ?", primaryRoleID).
			Find(&perms).Error; err != nil {
			return nil, err
		}
		permMap = make(map[int64][]domain.OrganizationMemberPagePermission)
		for _, p := range perms {
			permMap[p.RoleID] = append(permMap[p.RoleID], *p.ToDomain())
		}
	}

	result := make([]domain.UserOrganizationWithDetail, 0, len(rows))
	for _, row := range rows {
		if row.Organization == nil {
			continue
		}
		detail := domain.UserOrganizationWithDetail{
			UserOrganization: domain.UserOrganization{
				Organization: *row.Organization.ToDomain(),
				MemberID:     row.ID,
				RoleID:       row.RoleID,
				StatusID:     row.StatusID,
				IsOwner:      row.IsOwner,
				IsPrimary:    row.IsPrimary,
				JoinedAt:     row.JoinedAt,
			},
		}
		if row.Role != nil {
			detail.RoleName = row.Role.Name
		}
		if row.IsPrimary && permMap != nil {
			detail.PagePermissions = permMap[row.RoleID]
		}
		result = append(result, detail)
	}

	return result, nil
}

func (r *organizationRepository) FindPrimaryOrgWithDetails(ctx context.Context, userID int64) (*domain.PrimaryOrgPermissions, error) {
	var row models.OrganizationMemberModel

	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_primary = ?", userID, true).
		Preload("Organization").
		Preload("Role").
		First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	var perms []models.OrganizationMemberPagePermissionModel
	if err := r.db.WithContext(ctx).
		Where("role_id = ?", row.RoleID).
		Find(&perms).Error; err != nil {
		return nil, err
	}

	pagePermissions := make([]domain.OrganizationMemberPagePermission, len(perms))
	for i, p := range perms {
		pagePermissions[i] = *p.ToDomain()
	}

	result := &domain.PrimaryOrgPermissions{
		PagePermissions: pagePermissions,
	}
	if row.Organization != nil {
		result.OrganizationID = row.Organization.ID
	}
	if row.Role != nil {
		result.RoleName = row.Role.Name
	}

	return result, nil
}

func (r *organizationRepository) FindMembers(ctx context.Context, orgID int64, opts query.QueryOptions) ([]domain.OrganizationMember, int64, error) {
	var rows []models.OrganizationMemberModel
	var total int64

	base := r.db.WithContext(ctx).Model(&models.OrganizationMemberModel{}).Where("organization_id = ?", orgID)

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := gormq.ApplyToGorm(base, opts).Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	members := make([]domain.OrganizationMember, len(rows))
	for i, row := range rows {
		members[i] = *row.ToDomain()
	}

	return members, total, nil
}

func (r *organizationRepository) CreateMember(ctx context.Context, member *domain.OrganizationMember) error {
	now := time.Now()
	row := models.OrganizationMemberModel{
		OrganizationID: member.OrganizationID,
		UserID:         member.UserID,
		RoleID:         member.RoleID,
		StatusID:       member.StatusID,
		IsOwner:        member.IsOwner,
		IsPrimary:      member.IsPrimary,
		InvitedBy:      member.InvitedBy,
		InvitedAt:      &now,
	}
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return err
	}
	*member = *row.ToDomain()
	return nil
}

func (r *organizationRepository) UpdateMember(ctx context.Context, orgID int64, memberID int64, req *domain.UpdateMemberReq) error {
	return r.db.WithContext(ctx).Model(&models.OrganizationMemberModel{}).
		Where("id = ? AND organization_id = ?", memberID, orgID).
		Updates(utils.StructToMap(req)).Error
}

func (r *organizationRepository) DeleteMember(ctx context.Context, orgID int64, memberID int64) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND organization_id = ?", memberID, orgID).
		Delete(&models.OrganizationMemberModel{}).Error
}
