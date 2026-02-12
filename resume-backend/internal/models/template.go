package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// StringArray is a custom type for PostgreSQL text[] arrays
type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "{}", nil
	}
	// Format as PostgreSQL array literal: {"val1","val2"}
	var escaped []string
	for _, s := range a {
		// Escape quotes in strings
		s = strings.ReplaceAll(s, "\"", "\\\"")
		escaped = append(escaped, "\""+s+"\"")
	}
	return "{" + strings.Join(escaped, ",") + "}", nil
}

func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}

	var str string
	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	default:
		return errors.New("incompatible type for StringArray")
	}

	// Parse PostgreSQL array format: {val1,val2} or {"val1","val2"}
	str = strings.Trim(str, "{}")
	if str == "" {
		*a = []string{}
		return nil
	}

	// Simple parsing - handles basic cases
	var result []string
	for _, part := range strings.Split(str, ",") {
		part = strings.Trim(part, "\"")
		result = append(result, part)
	}
	*a = result
	return nil
}

type TemplateCategory string

const (
	CategoryModern     TemplateCategory = "modern"
	CategoryClassic    TemplateCategory = "classic"
	CategoryCreative   TemplateCategory = "creative"
	CategoryTech       TemplateCategory = "tech"
	CategoryExecutive  TemplateCategory = "executive"
	CategoryAcademic   TemplateCategory = "academic"
	CategoryMinimalist TemplateCategory = "minimalist"
)

type Template struct {
	ID              uuid.UUID        `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name            string           `gorm:"size:100;not null" json:"name"`
	Category        TemplateCategory `gorm:"type:varchar(50);not null" json:"category"`
	Description     string           `gorm:"size:500" json:"description"`
	PreviewImageURL string           `gorm:"size:500" json:"preview_image_url"`
	Config          *TemplateConfig  `gorm:"type:jsonb" json:"config"`
	IsPremium       bool             `gorm:"default:false" json:"is_premium"`
	ATSScore        int              `gorm:"default:80" json:"ats_score"`
	BestFor         StringArray      `gorm:"type:text[]" json:"best_for"`
	IsActive        bool             `gorm:"default:true" json:"is_active"`
	CreatedAt       time.Time        `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time        `gorm:"autoUpdateTime" json:"updated_at"`
}

func (t *Template) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// TemplateConfig defines the template appearance and layout
type TemplateConfig struct {
	Layout  LayoutConfig  `json:"layout"`
	Style   StyleConfig   `json:"style"`
	Section SectionConfig `json:"section"`
}

func (tc TemplateConfig) Value() (driver.Value, error) {
	return json.Marshal(tc)
}

func (tc *TemplateConfig) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, tc)
}

type LayoutConfig struct {
	Columns       int      `json:"columns"`        // 1 or 2
	PhotoPosition string   `json:"photo_position"` // left, right, top, none
	SectionOrder  []string `json:"section_order"`
	PageBreak     string   `json:"page_break"`     // auto, avoid, force
	LayoutVariant string   `json:"layout_variant"` // standard, sidebar_dark
}

type StyleConfig struct {
	Colors     ColorConfig     `json:"colors"`
	Typography TypographyConfig `json:"typography"`
	Spacing    SpacingConfig    `json:"spacing"`
	Decoration DecorationConfig `json:"decoration"`
}

type ColorConfig struct {
	Primary    string `json:"primary"`
	Secondary  string `json:"secondary"`
	Text       string `json:"text"`
	Background string `json:"background"`
	Accent     string `json:"accent"`
}

type TypographyConfig struct {
	HeadingFont   string  `json:"heading_font"`
	BodyFont      string  `json:"body_font"`
	BaseFontSize  int     `json:"base_font_size"`  // in pt
	LineHeight    float64 `json:"line_height"`
	LetterSpacing float64 `json:"letter_spacing"`
}

type SpacingConfig struct {
	Margins     string `json:"margins"`      // narrow, normal, wide
	SectionGap  int    `json:"section_gap"`  // in px
	ElementGap  int    `json:"element_gap"`  // in px
	PagePadding int    `json:"page_padding"` // in px
}

type DecorationConfig struct {
	Dividers    string `json:"dividers"`     // line, dots, none
	BulletStyle string `json:"bullet_style"` // disc, circle, square, dash
	UseIcons    bool   `json:"use_icons"`
	BorderStyle string `json:"border_style"` // none, thin, medium
}

type SectionConfig struct {
	Personal      SectionSettings `json:"personal"`
	Summary       SectionSettings `json:"summary"`
	Experience    SectionSettings `json:"experience"`
	Education     SectionSettings `json:"education"`
	Skills        SectionSettings `json:"skills"`
	Projects      SectionSettings `json:"projects"`
	Certifications SectionSettings `json:"certifications"`
}

type SectionSettings struct {
	Enabled    bool   `json:"enabled"`
	Label      string `json:"label"`
	Icon       string `json:"icon,omitempty"`
	LayoutType string `json:"layout_type"` // list, grid, timeline
}
