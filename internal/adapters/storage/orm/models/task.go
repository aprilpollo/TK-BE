package models

import (
	"aprilpollo/internal/core/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TasksModel struct {
	ID          int64     `gorm:"primaryKey"`
	ProjectID   int64     `gorm:"not null;index:idx_ticket_project"`
	Title       string    `gorm:"not null;size:500"`
	Description string    `gorm:"type:text"`
	Key         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	StatusID    int64     `gorm:"not null;index"`
	PriorityID  int64     `gorm:"not null;index"`
	ParentID    *int64    `gorm:"index"`
	Position    int       `gorm:"default:0;index"`
	StartDate   *int64    `gorm:"index"`
	EndDate     *int64    `gorm:"index"`
	AllDay      *bool     `gorm:"default:true"`

	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Project  *ProjectModel      `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	Status   *TaskStatusModel   `gorm:"foreignKey:StatusID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
	Priority *TaskPriorityModel `gorm:"foreignKey:PriorityID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
	Parent   *TasksModel        `gorm:"foreignKey:ParentID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
	Assigns  []TaskAssignModel  `gorm:"foreignKey:TaskID"`
}

type TaskAssignModel struct {
	ID        int64 `gorm:"primaryKey"`
	TaskID    int64 `gorm:"not null;index"`
	UserID    int64 `gorm:"not null;index"`
	StatusID  int64 `gorm:"not null;index;default:1"`
	IsLeader  bool  `gorm:"default:false"`
	JoinedAt  *time.Time
	InvitedAt *time.Time
	InvitedBy *int64

	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Task *TasksModel `gorm:"foreignKey:TaskID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	User *UserModel  `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
}

type TaskCommentModel struct {
	ID      int64  `gorm:"primaryKey"`
	TaskID  int64  `gorm:"not null;index:idx_comment_task"`
	UserID  int64  `gorm:"not null;index"`
	Content string `gorm:"not null;type:text"`

	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Task *TasksModel `gorm:"foreignKey:TaskID"`
	User *UserModel  `gorm:"foreignKey:UserID"`
}

type TaskAttachmentModel struct {
	ID         int64  `gorm:"primaryKey"`
	TaskID     int64  `gorm:"not null;index:idx_attachment_task"`
	Filename   string `gorm:"not null;size:255"`
	FilePath   string `gorm:"not null;size:500"`
	FileSize   int64  `gorm:"not null"`
	MimeType   string `gorm:"size:100"`
	UploadedBy int64  `gorm:"not null;index"`

	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Task     *TasksModel `gorm:"foreignKey:TaskID"`
	Uploader *UserModel  `gorm:"foreignKey:UploadedBy"`
}

type TaskStatusModel struct {
	ID          int64     `gorm:"primaryKey"`
	UUID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	ProjectID   int64     `gorm:"not null"`
	Name        string    `gorm:"not null;size:50;"`
	Description string    `gorm:"type:text;size:255"`
	Color       string    `gorm:"size:7;default:'#52525B'"`
	Position    int       `gorm:"default:0;index"`

	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Project *ProjectModel `gorm:"foreignKey:ProjectID"`
}

type TaskPriorityModel struct {
	ID          int64  `gorm:"primaryKey;autoIncrement:false"` // manually set IDs
	Name        string `gorm:"not null;size:50;uniqueIndex:idx_priority_project_name"`
	Description string `gorm:"type:text"`
	Color       string `gorm:"size:7;default:'#52525B'"`
}

func (TasksModel) TableName() string {
	return "tasks"
}

func (TaskAssignModel) TableName() string {
	return "task_assignments"
}

func (TaskCommentModel) TableName() string {
	return "task_comments"
}

func (TaskAttachmentModel) TableName() string {
	return "task_attachments"
}

func (TaskStatusModel) TableName() string {
	return "task_statuses"
}

func (TaskPriorityModel) TableName() string {
	return "task_priorities"
}

func (m *TaskStatusModel) ToDomain() *domain.TaskStatus {
	if m == nil {
		return nil
	}
	return &domain.TaskStatus{
		ID:          m.ID,
		UUID:        m.UUID,
		ProjectID:   m.ProjectID,
		Name:        m.Name,
		Description: m.Description,
		Color:       m.Color,
		Position:    m.Position,
	}
}

func (m *TaskPriorityModel) ToDomain() *domain.TaskPriority {
	if m == nil {
		return nil
	}
	return &domain.TaskPriority{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		Color:       m.Color,
	}
}

func (m *TasksModel) ToDomain() *domain.Task {
	if m == nil {
		return nil
	}

	status := domain.TaskStatus{ID: m.StatusID}
	if m.Status != nil {
		status = *m.Status.ToDomain()
	}

	priority := domain.TaskPriority{ID: m.PriorityID}
	if m.Priority != nil {
		priority = *m.Priority.ToDomain()
	}

	assigns := make([]domain.TaskAssign, 0, len(m.Assigns))
	for _, assign := range m.Assigns {
		item := domain.TaskAssign{ID: assign.UserID}

		if assign.User != nil {
			item.Name = assign.User.DisplayName
			item.Email = assign.User.Email
			if assign.User.Avatar != nil {
				item.Avatar = *assign.User.Avatar
			}
		}

		assigns = append(assigns, item)
	}

	return &domain.Task{
		ID:          m.ID,
		ProjectID:   m.ProjectID,
		Key:         m.Key,
		Title:       m.Title,
		Description: m.Description,
		StatusID:    m.StatusID,
		PriorityID:  m.PriorityID,
		ParentID:    m.ParentID,
		Position:    m.Position,
		StartDate:   m.StartDate,
		EndDate:     m.EndDate,
		AllDay:      m.AllDay != nil && *m.AllDay,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		Status:      status,
		Priority:    priority,
		Assigns:     assigns,
	}
}
