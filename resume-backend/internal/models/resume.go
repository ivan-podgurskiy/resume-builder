package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// JSONB type for PostgreSQL JSONB columns
type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, j)
}

type Resume struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Title       string         `gorm:"size:255;not null" json:"title"`
	IsMaster    bool           `gorm:"default:false" json:"is_master"`
	TemplateID  *uuid.UUID     `gorm:"type:uuid" json:"template_id,omitempty"`
	Data        *ResumeData    `gorm:"type:jsonb" json:"data"`
	StyleConfig JSONB          `gorm:"type:jsonb" json:"style_config,omitempty"`
	IsPublic    bool           `gorm:"default:false" json:"is_public"`
	PublicSlug  *string        `gorm:"size:100;uniqueIndex" json:"public_slug,omitempty"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User     *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Template *Template `gorm:"foreignKey:TemplateID" json:"template,omitempty"`
	Versions []ResumeVersion `gorm:"foreignKey:ResumeID" json:"versions,omitempty"`
}

func (r *Resume) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// ResumeData contains the full resume content structure
type ResumeData struct {
	PersonalInfo   PersonalInfo    `json:"personal_info"`
	Summary        string          `json:"summary"`
	Experience     []Experience    `json:"experience"`
	Education      []Education     `json:"education"`
	Skills         Skills          `json:"skills"`
	Certifications []Certification `json:"certifications,omitempty"`
	Projects       []Project       `json:"projects,omitempty"`
	Publications   []Publication   `json:"publications,omitempty"`
	Awards         []Award         `json:"awards,omitempty"`
	Volunteer      []Volunteer     `json:"volunteer,omitempty"`
	CustomSections []CustomSection `json:"custom_sections,omitempty"`
}

func (rd ResumeData) Value() (driver.Value, error) {
	return json.Marshal(rd)
}

func (rd *ResumeData) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, rd)
}

type PersonalInfo struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Title     string `json:"title"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Location  string `json:"location"`
	Website   string `json:"website,omitempty"`
	LinkedIn  string `json:"linkedin,omitempty"`
	GitHub    string `json:"github,omitempty"`
	PhotoURL  string `json:"photo_url,omitempty"`
}

type Experience struct {
	ID             string   `json:"id"`
	Company        string   `json:"company"`
	Position       string   `json:"position"`
	Location       string   `json:"location,omitempty"`
	StartDate      string   `json:"start_date"`
	EndDate        string   `json:"end_date,omitempty"`
	IsCurrent      bool     `json:"is_current"`
	Description    string   `json:"description"`
	Achievements   []string `json:"achievements,omitempty"`
	Technologies   []string `json:"technologies,omitempty"`
	Order          int      `json:"order"`
}

type Education struct {
	ID           string `json:"id"`
	Institution  string `json:"institution"`
	Degree       string `json:"degree"`
	FieldOfStudy string `json:"field_of_study"`
	Location     string `json:"location,omitempty"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date,omitempty"`
	IsCurrent    bool   `json:"is_current"`
	GPA          string `json:"gpa,omitempty"`
	Honors       string `json:"honors,omitempty"`
	Description  string `json:"description,omitempty"`
	Order        int    `json:"order"`
}

type Skills struct {
	Technical []Skill    `json:"technical"`
	Soft      []string   `json:"soft,omitempty"`
	Languages []Language `json:"languages,omitempty"`
}

type Skill struct {
	Name  string `json:"name"`
	Level string `json:"level,omitempty"` // beginner, intermediate, advanced, expert
}

type Language struct {
	Name        string `json:"name"`
	Proficiency string `json:"proficiency"` // native, fluent, advanced, intermediate, basic
}

type Certification struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Issuer       string `json:"issuer"`
	IssueDate    string `json:"issue_date"`
	ExpiryDate   string `json:"expiry_date,omitempty"`
	CredentialID string `json:"credential_id,omitempty"`
	URL          string `json:"url,omitempty"`
	Order        int    `json:"order"`
}

type Project struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	URL          string   `json:"url,omitempty"`
	RepoURL      string   `json:"repo_url,omitempty"`
	Technologies []string `json:"technologies,omitempty"`
	StartDate    string   `json:"start_date,omitempty"`
	EndDate      string   `json:"end_date,omitempty"`
	Highlights   []string `json:"highlights,omitempty"`
	Order        int      `json:"order"`
}

type Publication struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
	Date      string `json:"date"`
	URL       string `json:"url,omitempty"`
	Summary   string `json:"summary,omitempty"`
	Order     int    `json:"order"`
}

type Award struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Issuer      string `json:"issuer"`
	Date        string `json:"date"`
	Description string `json:"description,omitempty"`
	Order       int    `json:"order"`
}

type Volunteer struct {
	ID           string `json:"id"`
	Organization string `json:"organization"`
	Role         string `json:"role"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date,omitempty"`
	IsCurrent    bool   `json:"is_current"`
	Description  string `json:"description,omitempty"`
	Order        int    `json:"order"`
}

type CustomSection struct {
	ID      string             `json:"id"`
	Title   string             `json:"title"`
	Items   []CustomSectionItem `json:"items"`
	Order   int                `json:"order"`
}

type CustomSectionItem struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle,omitempty"`
	Date        string `json:"date,omitempty"`
	Description string `json:"description,omitempty"`
	Order       int    `json:"order"`
}

// ResumeVersion stores historical snapshots of resumes
type ResumeVersion struct {
	ID                uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ResumeID          uuid.UUID      `gorm:"type:uuid;not null;index" json:"resume_id"`
	VersionNumber     int            `gorm:"not null" json:"version_number"`
	DataSnapshot      *ResumeData    `gorm:"type:jsonb" json:"data_snapshot"`
	ChangeDescription string         `gorm:"size:500" json:"change_description,omitempty"`
	CreatedAt         time.Time      `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy         uuid.UUID      `gorm:"type:uuid" json:"created_by"`
}

func (rv *ResumeVersion) BeforeCreate(tx *gorm.DB) error {
	if rv.ID == uuid.Nil {
		rv.ID = uuid.New()
	}
	return nil
}
