package main

import (
	"aprilpollo/internal/adapters/config"
	"aprilpollo/internal/adapters/storage/orm"
	"aprilpollo/internal/adapters/storage/orm/models"
	"aprilpollo/internal/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"strconv"
    "time"
	"regexp"
	"os"
	"log"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// parseErrorMessage extracts the essential error information
func parseErrorMessage(err error) string {
	errorStr := err.Error()

	// Pattern to extract SQLSTATE code from any PostgreSQL error
	re := regexp.MustCompile(`\(SQLSTATE ([A-Z0-9]+)\)`)
	matches := re.FindStringSubmatch(errorStr)

	if len(matches) == 2 {
		return "SQLSTATE " + matches[1]
	}

	// If no SQLSTATE found, return the original error but truncated
	if len(errorStr) > 100 {
		return errorStr[:100] + "..."
	}

	return errorStr
}

// isDuplicateError checks if the error is a duplicate key violation (SQLSTATE 23505)
func isDuplicateError(err error) bool {
	return parseErrorMessage(err) == "SQLSTATE 23505"
}

// appendResultRow appends a row to the table based on the operation result
func appendResultRow(mTable table.Writer, name string, err error, failCount *int, successCount *int) {
	if err != nil {
		msg := parseErrorMessage(err)
		var status string
		if isDuplicateError(err) {
			status = text.Colors{text.FgYellow}.Sprint("✓ Skipped")
			*successCount++
		} else {
			status = text.Colors{text.FgRed}.Sprint("✗ Failed")
			*failCount++
		}
		mTable.AppendRow(table.Row{name, status, msg})
	} else {
		mTable.AppendRow(table.Row{name, text.Colors{text.FgGreen}.Sprint("✓ Created"), "SUCCESS"})
		*successCount++
	}
}

func main(){
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := orm.NewGormDB(cfg.Database, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}, true)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mTable := table.NewWriter()
	mTable.SetOutputMirror(os.Stdout)
	mTable.SetStyle(table.StyleRounded)
	mTable.Style().Title.Align = text.AlignCenter
	mTable.Style().Options.DoNotColorBordersAndSeparators = true
	mTable.Style().Options.DrawBorder = false
	mTable.Style().Options.SeparateColumns = true
	mTable.Style().Options.SeparateFooter = true
	mTable.Style().Options.SeparateHeader = true
	mTable.Style().Options.SeparateRows = false

	mTable.AppendHeader(table.Row{"NAME", "STATUS", "MESSAGE"})
	mTable.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, WidthMin: 20, AlignHeader: text.AlignLeft},
		{Number: 2, WidthMin: 20, AlignHeader: text.AlignLeft},
		{Number: 3, WidthMin: 20, AlignHeader: text.AlignLeft},
	})

	successCount := 0
	failCount := 0

	// ── Migration ────────────────────────────────────────────────────────────
	mTable.AppendRow(table.Row{"TABLES", "", "MIGRATION TABLES"})
	mTable.AppendRow(table.Row{"-", "-", "-"})

	for _, model := range models.All() {
		if err := db.Migrate(model); err != nil {
			mTable.AppendRow(table.Row{
				model.TableName(),
				text.Colors{text.FgRed}.Sprint("✗ Failed"),
				parseErrorMessage(err),
			})
			failCount++
		} else {
			mTable.AppendRow(table.Row{
				model.TableName(),
				text.Colors{text.FgGreen}.Sprint("✓ Migrated"),
				"SUCCESS",
			})
			successCount++
		}
	}

	mTable.AppendRow(table.Row{"", "", ""})
	mTable.AppendRow(table.Row{"INITIAL DATA", "", "INITIAL DATA"})
	mTable.AppendRow(table.Row{"-", "-", "-"})

	initDataDefault(db.GetDB(), mTable, &failCount, &successCount)

	mTable.AppendFooter(table.Row{"Summary", "", "Success: " + strconv.Itoa(successCount) + " Failed: " + strconv.Itoa(failCount)})
	mTable.Render()
}

func initDataDefault(db *gorm.DB, mTable table.Writer, failCount *int, successCount *int) {
	// MemberStatus
	memberStatuses := []models.OrganizationMemberStatusModel{
		{ID: 1, Name: "active", Description: "Active"},
		{ID: 2, Name: "invited", Description: "Invited"},
		{ID: 3, Name: "pending", Description: "Pending"},
		{ID: 4, Name: "inactive", Description: "Inactive"},
	}

	organizationMemberRoles := []models.OrganizationMemberRoleModel{
		{ID: 1, Name: "owner", Description: "Owner"},
		{ID: 2, Name: "admin", Description: "Admin"},
		{ID: 3, Name: "member", Description: "Member"},
	}

	projectStatuses := []models.ProjectStatusModel{
		{ID: 1, Name: "active", Description: "Active"},
		{ID: 2, Name: "inactive", Description: "Inactive"},
		{ID: 3, Name: "completed", Description: "Completed"},
		{ID: 4, Name: "cancelled", Description: "Cancelled"},
	}

	taskPriorities := []models.TaskPriorityModel{
		{ID: 1, Name: "low", Description: "Low", Color: "#52525B"},
		{ID: 2, Name: "medium", Description: "Medium", Color: "#12A150"},
		{ID: 3, Name: "high", Description: "High", Color: "#6020A0"},
	}

	userDefault := models.UserModel{
		Email:       "phonsing@gmail.com",
		FirstName:   "Phonsing",
		LastName:    "Taleman",
		DisplayName: "NO AH",
		IsActive:    true,
		IsVerified:  true,
	}

	// Insert Member Statuses
	for _, memberStatus := range memberStatuses {
		err := db.Create(&memberStatus).Error
		appendResultRow(mTable, "MEMBER STATUS ID: "+strconv.Itoa(int(memberStatus.ID)), err, failCount, successCount)
	}

	// Insert Organization Member Roles
	for _, role := range organizationMemberRoles {
		err := db.Create(&role).Error
		appendResultRow(mTable, "ORG MEMBER ROLE ID: "+strconv.Itoa(int(role.ID)), err, failCount, successCount)
	}

	// Page permissions per role
	// Pages: dashboard, members, projects, tasks
	// owner(1): full access | admin(2): limited | member(3): view only
	pagePermissions := []models.OrganizationMemberPagePermissionModel{
		// owner — full access
		{ID: 1, PageID: "dashboard", RoleID: 1, IsView: true, IsEdit: true, IsDelete: true},
		{ID: 2, PageID: "members", RoleID: 1, IsView: true, IsEdit: true, IsDelete: true},
		{ID: 3, PageID: "projects", RoleID: 1, IsView: true, IsEdit: true, IsDelete: true},
		{ID: 4, PageID: "tasks", RoleID: 1, IsView: true, IsEdit: true, IsDelete: true},
		// admin — no member management delete
		{ID: 5, PageID: "dashboard", RoleID: 2, IsView: true, IsEdit: true, IsDelete: false},
		{ID: 6, PageID: "members", RoleID: 2, IsView: true, IsEdit: false, IsDelete: false},
		{ID: 7, PageID: "projects", RoleID: 2, IsView: true, IsEdit: true, IsDelete: true},
		{ID: 8, PageID: "tasks", RoleID: 2, IsView: true, IsEdit: true, IsDelete: true},
		// member — view only, can edit own tasks
		{ID: 9, PageID: "dashboard", RoleID: 3, IsView: true, IsEdit: false, IsDelete: false},
		{ID: 10, PageID: "members", RoleID: 3, IsView: true, IsEdit: false, IsDelete: false},
		{ID: 11, PageID: "projects", RoleID: 3, IsView: true, IsEdit: false, IsDelete: false},
		{ID: 12, PageID: "tasks", RoleID: 3, IsView: true, IsEdit: true, IsDelete: false},
	}

	for _, pp := range pagePermissions {
		err := db.Create(&pp).Error
		appendResultRow(mTable, "PAGE PERMISSION ID: "+strconv.Itoa(int(pp.ID)), err, failCount, successCount)
	}

	// Insert Project Statuses (fixed variable shadowing)
	for _, ps := range projectStatuses {
		err := db.Create(&ps).Error
		appendResultRow(mTable, "PROJECT STATUS ID: "+strconv.Itoa(int(ps.ID)), err, failCount, successCount)
	}

	// Insert Task Priorities
	for _, tp := range taskPriorities {
		err := db.Create(&tp).Error
		appendResultRow(mTable, "TASK PRIORITY ID: "+strconv.Itoa(int(tp.ID)), err, failCount, successCount)
	}

	// Insert Default User
	if err := db.Create(&userDefault).Error; err != nil {
		appendResultRow(mTable, "DEFAULT USER", err, failCount, successCount)
		// If user already exists (duplicate) or any other error, skip Oauth/Organization/OrganizationMember
		return
	}
	appendResultRow(mTable, "DEFAULT USER", nil, failCount, successCount)

	// Create OAuth for default user
	hashPassword, err := utils.HashPassword("@aplps9921")
	if err != nil {
		mTable.AppendRow(table.Row{"PASSWORD HASH", text.Colors{text.FgRed}.Sprint("✗ Failed"), err.Error()})
		*failCount++
		return
	}

	oauthDefault := models.OauthModel{
		UserID:     userDefault.ID,
		Provider:   "basic",
		ProviderID: "1234567890",
		Email:      "phonsing@gmail.com",
		Password:   &hashPassword,
	}

	if err := db.Create(&oauthDefault).Error; err != nil {
		appendResultRow(mTable, "DEFAULT OAUTH", err, failCount, successCount)
		if !isDuplicateError(err) {
			return
		}
	} else {
		appendResultRow(mTable, "DEFAULT OAUTH", nil, failCount, successCount)
	}

	// Insert Default Organizations
	defaultOrganizations := []models.OrganizationModel{
		{
			Name:         "Default Organization 01",
			Description:  "Default Organization",
			Slug:         "default-organization-01",
			IsActive:     true,
			ContactEmail: "phonsing@gmail.com",
		},
		{
			Name:         "Default Organization 02",
			Description:  "Default Organization",
			Slug:         "default-organization-02",
			IsActive:     true,
			ContactEmail: "phonsing@gmail.com",
		},
		{
			Name:         "Default Organization 03",
			Description:  "Default Organization",
			Slug:         "default-organization-03",
			IsActive:     true,
			ContactEmail: "phonsing@gmail.com",
		},
	}

	now := time.Now()
	for index, org := range defaultOrganizations {
		if err := db.Create(&org).Error; err != nil {
			appendResultRow(mTable, "ORGANIZATION: "+org.Slug, err, failCount, successCount)
			if !isDuplicateError(err) {
				continue
			}
			// Fetch existing organization to get its ID
			if fetchErr := db.Where("slug = ?", org.Slug).First(&org).Error; fetchErr != nil {
				continue
			}
		} else {
			appendResultRow(mTable, "ORGANIZATION: "+org.Slug, nil, failCount, successCount)
		}

		isPrimary := index == 0

		defaultOrganizationMember := models.OrganizationMemberModel{
			OrganizationID: org.ID,
			UserID:         userDefault.ID,
			RoleID:         1,
			StatusID:       1,
			IsOwner:        true,
			IsPrimary:      isPrimary,
			InvitedAt:      &now,
			JoinedAt:       &now,
		}

		if err := db.Create(&defaultOrganizationMember).Error; err != nil {
			appendResultRow(mTable, "ORG MEMBER: "+org.Slug, err, failCount, successCount)
		} else {
			appendResultRow(mTable, "ORG MEMBER: "+org.Slug, nil, failCount, successCount)
		}
	}
}
