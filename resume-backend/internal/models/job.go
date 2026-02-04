package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JobType string

const (
	JobTypePDFExport  JobType = "pdf_export"
	JobTypeDOCXExport JobType = "docx_export"
	JobTypeAIExtract  JobType = "ai_extract"
	JobTypeAIImprove  JobType = "ai_improve"
	JobTypeEmail      JobType = "email"
)

type JobStatus string

const (
	JobStatusPending    JobStatus = "pending"
	JobStatusProcessing JobStatus = "processing"
	JobStatusCompleted  JobStatus = "completed"
	JobStatusFailed     JobStatus = "failed"
)

type Job struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Type        JobType        `gorm:"type:varchar(50);not null" json:"type"`
	Status      JobStatus      `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	Payload     JSONB          `gorm:"type:jsonb" json:"payload,omitempty"`
	Result      JSONB          `gorm:"type:jsonb" json:"result,omitempty"`
	Error       string         `gorm:"type:text" json:"error,omitempty"`
	Progress    int            `gorm:"default:0" json:"progress"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	CompletedAt *time.Time     `json:"completed_at,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (j *Job) BeforeCreate(tx *gorm.DB) error {
	if j.ID == uuid.Nil {
		j.ID = uuid.New()
	}
	return nil
}

// Export represents a completed export file
type Export struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ResumeID  uuid.UUID      `gorm:"type:uuid;not null;index" json:"resume_id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Format    string         `gorm:"type:varchar(20);not null" json:"format"` // pdf, docx, txt, html, json
	FileURL   string         `gorm:"size:500;not null" json:"file_url"`
	FileSize  int64          `json:"file_size"`
	FileHash  string         `gorm:"size:64" json:"file_hash"`
	Settings  JSONB          `gorm:"type:jsonb" json:"settings,omitempty"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	ExpiresAt *time.Time     `json:"expires_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Resume *Resume `gorm:"foreignKey:ResumeID" json:"resume,omitempty"`
	User   *User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (e *Export) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

// AIRequest tracks AI usage for analytics and rate limiting
type AIRequest struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID       uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	RequestType  string    `gorm:"type:varchar(50);not null" json:"request_type"`
	InputTokens  int       `json:"input_tokens"`
	OutputTokens int       `json:"output_tokens"`
	CostUSD      float64   `json:"cost_usd"`
	DurationMS   int       `json:"duration_ms"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relations
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (ar *AIRequest) BeforeCreate(tx *gorm.DB) error {
	if ar.ID == uuid.Nil {
		ar.ID = uuid.New()
	}
	return nil
}

// AuditLog tracks user actions for compliance
type AuditLog struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID     uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	Action     string    `gorm:"type:varchar(100);not null" json:"action"`
	EntityType string    `gorm:"type:varchar(50)" json:"entity_type"`
	EntityID   string    `gorm:"type:varchar(50)" json:"entity_id"`
	Metadata   JSONB     `gorm:"type:jsonb" json:"metadata,omitempty"`
	IPAddress  string    `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent  string    `gorm:"type:text" json:"user_agent"`
	CreatedAt  time.Time `gorm:"autoCreateTime;index" json:"created_at"`
}

func (al *AuditLog) BeforeCreate(tx *gorm.DB) error {
	if al.ID == uuid.Nil {
		al.ID = uuid.New()
	}
	return nil
}
