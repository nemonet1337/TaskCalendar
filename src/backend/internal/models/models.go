package models

import (
	"time"

	"gorm.io/gorm"
)

// User モデル
type User struct {
	ID        string `json:"id" gorm:"primaryKey;type:varchar(25)"`
	Email     string `json:"email" gorm:"unique;not null"`
	Username  string `json:"username" gorm:"unique;not null"`
	Password  string `json:"-" gorm:"not null"`
	FirstName string `json:"firstName" gorm:"not null"`
	LastName  string `json:"lastName" gorm:"not null"`
	Avatar    string `json:"avatar"`
	Role      UserRole `json:"role" gorm:"default:'MEMBER'"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	// Relations
	TeamMemberships []TeamMember `json:"teamMemberships" gorm:"foreignKey:UserID"`
	CreatedTeams    []Team       `json:"createdTeams" gorm:"foreignKey:CreatorID"`
	AssignedTasks   []Task       `json:"assignedTasks" gorm:"foreignKey:AssigneeID"`
	CreatedTasks    []Task       `json:"createdTasks" gorm:"foreignKey:CreatorID"`
	Events          []Event      `json:"events" gorm:"foreignKey:CreatorID"`
	Comments        []Comment    `json:"comments" gorm:"foreignKey:AuthorID"`
}

type UserRole string

const (
	UserRoleAdmin   UserRole = "ADMIN"
	UserRoleManager UserRole = "MANAGER"
	UserRoleMember  UserRole = "MEMBER"
)

// Team モデル
type Team struct {
	ID          string `json:"id" gorm:"primaryKey;type:varchar(25)"`
	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	CreatorID   string `json:"creatorId" gorm:"not null"`

	// Relations
	Creator User         `json:"creator" gorm:"foreignKey:CreatorID"`
	Members []TeamMember `json:"members" gorm:"foreignKey:TeamID"`
	Tasks   []Task       `json:"tasks" gorm:"foreignKey:TeamID"`
	Events  []Event      `json:"events" gorm:"foreignKey:TeamID"`
}

// TeamMember モデル
type TeamMember struct {
	ID       string           `json:"id" gorm:"primaryKey;type:varchar(25)"`
	UserID   string           `json:"userId" gorm:"not null"`
	TeamID   string           `json:"teamId" gorm:"not null"`
	Role     TeamMemberRole   `json:"role" gorm:"default:'MEMBER'"`
	Status   TeamMemberStatus `json:"status" gorm:"default:'ACTIVE'"`
	JoinedAt time.Time        `json:"joinedAt"`

	// Relations
	User User `json:"user" gorm:"foreignKey:UserID"`
	Team Team `json:"team" gorm:"foreignKey:TeamID"`
}

type TeamMemberRole string

const (
	TeamMemberRoleOwner  TeamMemberRole = "OWNER"
	TeamMemberRoleAdmin  TeamMemberRole = "ADMIN"
	TeamMemberRoleMember TeamMemberRole = "MEMBER"
)

type TeamMemberStatus string

const (
	TeamMemberStatusActive   TeamMemberStatus = "ACTIVE"
	TeamMemberStatusInactive TeamMemberStatus = "INACTIVE"
	TeamMemberStatusPending  TeamMemberStatus = "PENDING"
)

// Task モデル
type Task struct {
	ID          string `json:"id" gorm:"primaryKey;type:varchar(25)"`
	Title       string `json:"title" gorm:"not null"`
	Description string `json:"description"`
	Status      TaskStatus `json:"status" gorm:"default:'TODO'"`
	Priority    Priority `json:"priority" gorm:"default:'MEDIUM'"`
	DueDate     *time.Time `json:"dueDate"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	TeamID      string `json:"teamId" gorm:"not null"`
	CreatorID   string `json:"creatorId" gorm:"not null"`
	AssigneeID  *string `json:"assigneeId"`

	// Relations
	Team     Team      `json:"team" gorm:"foreignKey:TeamID"`
	Creator  User      `json:"creator" gorm:"foreignKey:CreatorID"`
	Assignee *User     `json:"assignee" gorm:"foreignKey:AssigneeID"`
	Comments []Comment `json:"comments" gorm:"foreignKey:TaskID"`
}

type TaskStatus string

const (
	TaskStatusTodo       TaskStatus = "TODO"
	TaskStatusInProgress TaskStatus = "IN_PROGRESS"
	TaskStatusInReview   TaskStatus = "IN_REVIEW"
	TaskStatusDone       TaskStatus = "DONE"
	TaskStatusCancelled  TaskStatus = "CANCELLED"
)

type Priority string

const (
	PriorityLow    Priority = "LOW"
	PriorityMedium Priority = "MEDIUM"
	PriorityHigh   Priority = "HIGH"
	PriorityUrgent Priority = "URGENT"
)

// Event モデル
type Event struct {
	ID          string `json:"id" gorm:"primaryKey;type:varchar(25)"`
	Title       string `json:"title" gorm:"not null"`
	Description string `json:"description"`
	StartDate   time.Time `json:"startDate" gorm:"not null"`
	EndDate     time.Time `json:"endDate" gorm:"not null"`
	IsRecurring bool   `json:"isRecurring" gorm:"default:false"`
	Recurrence  string `json:"recurrence"`
	Type        EventType `json:"type" gorm:"default:'MEETING'"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	TeamID      *string `json:"teamId"`
	CreatorID   string `json:"creatorId" gorm:"not null"`

	// Relations
	Team    *Team `json:"team" gorm:"foreignKey:TeamID"`
	Creator User  `json:"creator" gorm:"foreignKey:CreatorID"`
}

type EventType string

const (
	EventTypeMeeting  EventType = "MEETING"
	EventTypeDeadline EventType = "DEADLINE"
	EventTypeReminder EventType = "REMINDER"
	EventTypePersonal EventType = "PERSONAL"
)

// Comment モデル
type Comment struct {
	ID        string `json:"id" gorm:"primaryKey;type:varchar(25)"`
	Content   string `json:"content" gorm:"not null"`
	CreatedAt time.Time `json:"createdAt"`
	TaskID    string `json:"taskId" gorm:"not null"`
	AuthorID  string `json:"authorId" gorm:"not null"`

	// Relations
	Task   Task `json:"task" gorm:"foreignKey:TaskID"`
	Author User `json:"author" gorm:"foreignKey:AuthorID"`
}

// BeforeCreate フック - ID生成
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = generateID()
	}
	return nil
}

func (t *Team) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = generateID()
	}
	return nil
}

func (tm *TeamMember) BeforeCreate(tx *gorm.DB) error {
	if tm.ID == "" {
		tm.ID = generateID()
	}
	return nil
}

func (t *Task) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = generateID()
	}
	return nil
}

func (e *Event) BeforeCreate(tx *gorm.DB) error {
	if e.ID == "" {
		e.ID = generateID()
	}
	return nil
}

func (c *Comment) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = generateID()
	}
	return nil
}

// 簡単なID生成関数（実際のプロダクションではより堅牢な実装を推奨）
func generateID() string {
	// 実装は省略 - 実際にはnanoid等を使用
	return "temp_id"
}
