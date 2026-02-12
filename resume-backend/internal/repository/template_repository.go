package repository

import (
	"github.com/google/uuid"
	"github.com/resume-builder/backend/internal/models"
	"gorm.io/gorm"
)

type TemplateRepository struct {
	db *gorm.DB
}

func NewTemplateRepository(db *gorm.DB) *TemplateRepository {
	return &TemplateRepository{db: db}
}

func (r *TemplateRepository) Create(template *models.Template) error {
	return r.db.Create(template).Error
}

func (r *TemplateRepository) FindByID(id uuid.UUID) (*models.Template, error) {
	var template models.Template
	err := r.db.First(&template, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *TemplateRepository) FindAll(includeInactive bool) ([]models.Template, error) {
	var templates []models.Template
	query := r.db.Model(&models.Template{})
	if !includeInactive {
		query = query.Where("is_active = ?", true)
	}
	err := query.Order("is_premium ASC, name ASC").Find(&templates).Error
	return templates, err
}

func (r *TemplateRepository) FindByCategory(category models.TemplateCategory) ([]models.Template, error) {
	var templates []models.Template
	err := r.db.Where("category = ? AND is_active = ?", category, true).
		Order("is_premium ASC, name ASC").
		Find(&templates).Error
	return templates, err
}

func (r *TemplateRepository) FindFreeTemplates() ([]models.Template, error) {
	var templates []models.Template
	err := r.db.Where("is_premium = ? AND is_active = ?", false, true).
		Order("name ASC").
		Find(&templates).Error
	return templates, err
}

func (r *TemplateRepository) Update(template *models.Template) error {
	return r.db.Save(template).Error
}

func (r *TemplateRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Template{}, "id = ?", id).Error
}

// SeedDefaultTemplates creates the initial set of templates
func (r *TemplateRepository) SeedDefaultTemplates() error {
	templates := []models.Template{
		{
			Name:        "Modern Minimalist",
			Category:    models.CategoryModern,
			Description: "Single column, clean typography, subtle accents. ATS-friendly design perfect for any industry.",
			IsPremium:   false,
			ATSScore:    95,
			BestFor:     models.StringArray{"Software Engineer", "Product Manager", "Designer"},
			Config: &models.TemplateConfig{
				Layout: models.LayoutConfig{
					Columns:       1,
					PhotoPosition: "none",
					SectionOrder:  []string{"personal", "summary", "experience", "education", "skills", "projects"},
					PageBreak:     "auto",
				},
				Style: models.StyleConfig{
					Colors: models.ColorConfig{
						Primary:    "#2563eb",
						Secondary:  "#1e40af",
						Text:       "#111827",
						Background: "#ffffff",
						Accent:     "#60a5fa",
					},
					Typography: models.TypographyConfig{
						HeadingFont:   "Inter",
						BodyFont:      "Inter",
						BaseFontSize:  11,
						LineHeight:    1.5,
						LetterSpacing: 0,
					},
					Spacing: models.SpacingConfig{
						Margins:     "normal",
						SectionGap:  16,
						ElementGap:  8,
						PagePadding: 40,
					},
					Decoration: models.DecorationConfig{
						Dividers:    "line",
						BulletStyle: "disc",
						UseIcons:    false,
						BorderStyle: "none",
					},
				},
				Section: models.SectionConfig{
					Personal:       models.SectionSettings{Enabled: true, Label: "Contact", LayoutType: "list"},
					Summary:        models.SectionSettings{Enabled: true, Label: "Professional Summary", LayoutType: "list"},
					Experience:     models.SectionSettings{Enabled: true, Label: "Experience", LayoutType: "timeline"},
					Education:      models.SectionSettings{Enabled: true, Label: "Education", LayoutType: "list"},
					Skills:         models.SectionSettings{Enabled: true, Label: "Skills", LayoutType: "grid"},
					Projects:       models.SectionSettings{Enabled: true, Label: "Projects", LayoutType: "list"},
					Certifications: models.SectionSettings{Enabled: true, Label: "Certifications", LayoutType: "list"},
				},
			},
		},
		{
			Name:        "Professional Classic",
			Category:    models.CategoryClassic,
			Description: "Two column layout with traditional formatting. Perfect for corporate environments.",
			IsPremium:   false,
			ATSScore:    90,
			BestFor:     models.StringArray{"Finance", "Consulting", "Legal", "Management"},
			Config: &models.TemplateConfig{
				Layout: models.LayoutConfig{
					Columns:       2,
					PhotoPosition: "left",
					SectionOrder:  []string{"personal", "summary", "experience", "education", "skills", "certifications"},
					PageBreak:     "auto",
				},
				Style: models.StyleConfig{
					Colors: models.ColorConfig{
						Primary:    "#1e3a5f",
						Secondary:  "#2d5a87",
						Text:       "#333333",
						Background: "#ffffff",
						Accent:     "#4a90d9",
					},
					Typography: models.TypographyConfig{
						HeadingFont:   "Georgia",
						BodyFont:      "Inter",
						BaseFontSize:  10,
						LineHeight:    1.4,
						LetterSpacing: 0,
					},
					Spacing: models.SpacingConfig{
						Margins:     "normal",
						SectionGap:  12,
						ElementGap:  6,
						PagePadding: 36,
					},
					Decoration: models.DecorationConfig{
						Dividers:    "line",
						BulletStyle: "disc",
						UseIcons:    false,
						BorderStyle: "thin",
					},
				},
				Section: models.SectionConfig{
					Personal:       models.SectionSettings{Enabled: true, Label: "Contact Information", LayoutType: "list"},
					Summary:        models.SectionSettings{Enabled: true, Label: "Professional Summary", LayoutType: "list"},
					Experience:     models.SectionSettings{Enabled: true, Label: "Professional Experience", LayoutType: "list"},
					Education:      models.SectionSettings{Enabled: true, Label: "Education", LayoutType: "list"},
					Skills:         models.SectionSettings{Enabled: true, Label: "Skills", LayoutType: "list"},
					Projects:       models.SectionSettings{Enabled: false, Label: "Projects", LayoutType: "list"},
					Certifications: models.SectionSettings{Enabled: true, Label: "Certifications", LayoutType: "list"},
				},
			},
		},
		{
			Name:        "Tech Optimized",
			Category:    models.CategoryTech,
			Description: "Skills matrix, project showcase, and GitHub integration. Designed for technical roles.",
			IsPremium:   false,
			ATSScore:    92,
			BestFor:     models.StringArray{"Software Developer", "DevOps", "Data Scientist", "Engineer"},
			Config: &models.TemplateConfig{
				Layout: models.LayoutConfig{
					Columns:       1,
					PhotoPosition: "none",
					SectionOrder:  []string{"personal", "summary", "skills", "experience", "projects", "education", "certifications"},
					PageBreak:     "auto",
				},
				Style: models.StyleConfig{
					Colors: models.ColorConfig{
						Primary:    "#0d9488",
						Secondary:  "#115e59",
						Text:       "#1f2937",
						Background: "#ffffff",
						Accent:     "#14b8a6",
					},
					Typography: models.TypographyConfig{
						HeadingFont:   "JetBrains Mono",
						BodyFont:      "Inter",
						BaseFontSize:  10,
						LineHeight:    1.5,
						LetterSpacing: 0,
					},
					Spacing: models.SpacingConfig{
						Margins:     "normal",
						SectionGap:  14,
						ElementGap:  8,
						PagePadding: 40,
					},
					Decoration: models.DecorationConfig{
						Dividers:    "line",
						BulletStyle: "dash",
						UseIcons:    true,
						BorderStyle: "none",
					},
				},
				Section: models.SectionConfig{
					Personal:       models.SectionSettings{Enabled: true, Label: "Contact", Icon: "user", LayoutType: "list"},
					Summary:        models.SectionSettings{Enabled: true, Label: "About", Icon: "info", LayoutType: "list"},
					Experience:     models.SectionSettings{Enabled: true, Label: "Experience", Icon: "briefcase", LayoutType: "timeline"},
					Education:      models.SectionSettings{Enabled: true, Label: "Education", Icon: "graduation-cap", LayoutType: "list"},
					Skills:         models.SectionSettings{Enabled: true, Label: "Technical Skills", Icon: "code", LayoutType: "grid"},
					Projects:       models.SectionSettings{Enabled: true, Label: "Projects", Icon: "folder", LayoutType: "grid"},
					Certifications: models.SectionSettings{Enabled: true, Label: "Certifications", Icon: "award", LayoutType: "list"},
				},
			},
		},
		// Modern
		{
			Name:        "Modern Professional",
			Category:    models.CategoryModern,
			Description: "Clean two-column layout with strong typography. Ideal for mid-career professionals.",
			IsPremium:   false,
			ATSScore:    93,
			BestFor:     models.StringArray{"Product Manager", "Marketing", "Consultant"},
			Config: &models.TemplateConfig{
				Layout: models.LayoutConfig{
					Columns:       2,
					PhotoPosition: "right",
					SectionOrder:  []string{"personal", "summary", "experience", "education", "skills", "certifications"},
					PageBreak:     "auto",
				},
				Style: models.StyleConfig{
					Colors: models.ColorConfig{
						Primary:    "#7c3aed",
						Secondary:  "#5b21b6",
						Text:       "#1e293b",
						Background: "#ffffff",
						Accent:     "#a78bfa",
					},
					Typography: models.TypographyConfig{
						HeadingFont:   "Inter",
						BodyFont:      "Inter",
						BaseFontSize:  10,
						LineHeight:    1.45,
						LetterSpacing: 0,
					},
					Spacing: models.SpacingConfig{
						Margins:     "normal",
						SectionGap:  14,
						ElementGap:  6,
						PagePadding: 36,
					},
					Decoration: models.DecorationConfig{
						Dividers:    "line",
						BulletStyle: "disc",
						UseIcons:    false,
						BorderStyle: "thin",
					},
				},
				Section: models.SectionConfig{
					Personal:       models.SectionSettings{Enabled: true, Label: "Contact", LayoutType: "list"},
					Summary:        models.SectionSettings{Enabled: true, Label: "Profile", LayoutType: "list"},
					Experience:     models.SectionSettings{Enabled: true, Label: "Experience", LayoutType: "list"},
					Education:      models.SectionSettings{Enabled: true, Label: "Education", LayoutType: "list"},
					Skills:         models.SectionSettings{Enabled: true, Label: "Skills", LayoutType: "list"},
					Projects:       models.SectionSettings{Enabled: false, Label: "Projects", LayoutType: "list"},
					Certifications: models.SectionSettings{Enabled: true, Label: "Certifications", LayoutType: "list"},
				},
			},
		},
		// Classic
		{
			Name:        "Classic Elegant",
			Category:    models.CategoryClassic,
			Description: "Traditional single-column format with serif headings. Timeless for established professionals.",
			IsPremium:   false,
			ATSScore:    88,
			BestFor:     models.StringArray{"Legal", "Finance", "Healthcare"},
			Config: &models.TemplateConfig{
				Layout: models.LayoutConfig{
					Columns:       1,
					PhotoPosition: "none",
					SectionOrder:  []string{"personal", "summary", "experience", "education", "skills", "certifications"},
					PageBreak:     "auto",
				},
				Style: models.StyleConfig{
					Colors: models.ColorConfig{
						Primary:    "#1e293b",
						Secondary:  "#334155",
						Text:       "#334155",
						Background: "#ffffff",
						Accent:     "#64748b",
					},
					Typography: models.TypographyConfig{
						HeadingFont:   "Georgia",
						BodyFont:      "Georgia",
						BaseFontSize:  11,
						LineHeight:    1.4,
						LetterSpacing: 0,
					},
					Spacing: models.SpacingConfig{
						Margins:     "wide",
						SectionGap:  18,
						ElementGap:  8,
						PagePadding: 48,
					},
					Decoration: models.DecorationConfig{
						Dividers:    "line",
						BulletStyle: "disc",
						UseIcons:    false,
						BorderStyle: "thin",
					},
				},
				Section: models.SectionConfig{
					Personal:       models.SectionSettings{Enabled: true, Label: "Contact Information", LayoutType: "list"},
					Summary:        models.SectionSettings{Enabled: true, Label: "Professional Summary", LayoutType: "list"},
					Experience:     models.SectionSettings{Enabled: true, Label: "Work Experience", LayoutType: "list"},
					Education:      models.SectionSettings{Enabled: true, Label: "Education", LayoutType: "list"},
					Skills:         models.SectionSettings{Enabled: true, Label: "Skills & Expertise", LayoutType: "list"},
					Projects:       models.SectionSettings{Enabled: false, Label: "Projects", LayoutType: "list"},
					Certifications: models.SectionSettings{Enabled: true, Label: "Certifications", LayoutType: "list"},
				},
			},
		},
		// Creative (2)
		{
			Name:        "Creative Bold",
			Category:    models.CategoryCreative,
			Description: "Vibrant accent colors and contemporary layout. Stand out in creative industries.",
			IsPremium:   false,
			ATSScore:    85,
			BestFor:     models.StringArray{"Designer", "Creative Director", "Artist"},
			Config: &models.TemplateConfig{
				Layout: models.LayoutConfig{
					Columns:       2,
					PhotoPosition: "left",
					SectionOrder:  []string{"personal", "summary", "experience", "skills", "education"},
					PageBreak:     "auto",
				},
				Style: models.StyleConfig{
					Colors: models.ColorConfig{
						Primary:    "#dc2626",
						Secondary:  "#b91c1c",
						Text:       "#1f2937",
						Background: "#ffffff",
						Accent:     "#f87171",
					},
					Typography: models.TypographyConfig{
						HeadingFont:   "Inter",
						BodyFont:      "Inter",
						BaseFontSize:  10,
						LineHeight:    1.5,
						LetterSpacing: 0,
					},
					Spacing: models.SpacingConfig{
						Margins:     "normal",
						SectionGap:  16,
						ElementGap:  8,
						PagePadding: 40,
					},
					Decoration: models.DecorationConfig{
						Dividers:    "line",
						BulletStyle: "square",
						UseIcons:    false,
						BorderStyle: "medium",
					},
				},
				Section: models.SectionConfig{
					Personal:   models.SectionSettings{Enabled: true, Label: "Contact", LayoutType: "list"},
					Summary:    models.SectionSettings{Enabled: true, Label: "About", LayoutType: "list"},
					Experience: models.SectionSettings{Enabled: true, Label: "Experience", LayoutType: "list"},
					Education:  models.SectionSettings{Enabled: true, Label: "Education", LayoutType: "list"},
					Skills:     models.SectionSettings{Enabled: true, Label: "Skills", LayoutType: "grid"},
					Projects:   models.SectionSettings{Enabled: false, Label: "Projects", LayoutType: "list"},
					Certifications: models.SectionSettings{Enabled: false, Label: "Certifications", LayoutType: "list"},
				},
			},
		},
		{
			Name:        "Creative Modern",
			Category:    models.CategoryCreative,
			Description: "Fresh design with teal accents. Perfect for design and marketing roles.",
			IsPremium:   false,
			ATSScore:    87,
			BestFor:     models.StringArray{"UX Designer", "Content Strategist", "Brand Manager"},
			Config: &models.TemplateConfig{
				Layout: models.LayoutConfig{
					Columns:       1,
					PhotoPosition: "top",
					SectionOrder:  []string{"personal", "summary", "experience", "skills", "education", "projects"},
					PageBreak:     "auto",
				},
				Style: models.StyleConfig{
					Colors: models.ColorConfig{
						Primary:    "#0f766e",
						Secondary:  "#134e4a",
						Text:       "#1e293b",
						Background: "#f8fafc",
						Accent:     "#2dd4bf",
					},
					Typography: models.TypographyConfig{
						HeadingFont:   "Inter",
						BodyFont:      "Inter",
						BaseFontSize:  10,
						LineHeight:    1.55,
						LetterSpacing: 0,
					},
					Spacing: models.SpacingConfig{
						Margins:     "normal",
						SectionGap:  18,
						ElementGap:  10,
						PagePadding: 44,
					},
					Decoration: models.DecorationConfig{
						Dividers:    "dots",
						BulletStyle: "circle",
						UseIcons:    false,
						BorderStyle: "none",
					},
				},
				Section: models.SectionConfig{
					Personal:       models.SectionSettings{Enabled: true, Label: "Contact", LayoutType: "list"},
					Summary:        models.SectionSettings{Enabled: true, Label: "Summary", LayoutType: "list"},
					Experience:     models.SectionSettings{Enabled: true, Label: "Experience", LayoutType: "timeline"},
					Education:      models.SectionSettings{Enabled: true, Label: "Education", LayoutType: "list"},
					Skills:         models.SectionSettings{Enabled: true, Label: "Skills", LayoutType: "grid"},
					Projects:       models.SectionSettings{Enabled: true, Label: "Projects", LayoutType: "list"},
					Certifications: models.SectionSettings{Enabled: false, Label: "Certifications", LayoutType: "list"},
				},
			},
		},
		// Minimalist (2)
		{
			Name:        "Clean Minimalist",
			Category:    models.CategoryMinimalist,
			Description: "Ultra-clean single column. Maximum readability, minimal decoration.",
			IsPremium:   false,
			ATSScore:    96,
			BestFor:     models.StringArray{"Any Role", "ATS-Optimized"},
			Config: &models.TemplateConfig{
				Layout: models.LayoutConfig{
					Columns:       1,
					PhotoPosition: "none",
					SectionOrder:  []string{"personal", "summary", "experience", "education", "skills"},
					PageBreak:     "auto",
				},
				Style: models.StyleConfig{
					Colors: models.ColorConfig{
						Primary:    "#18181b",
						Secondary:  "#27272a",
						Text:       "#3f3f46",
						Background: "#ffffff",
						Accent:     "#71717a",
					},
					Typography: models.TypographyConfig{
						HeadingFont:   "Inter",
						BodyFont:      "Inter",
						BaseFontSize:  11,
						LineHeight:    1.6,
						LetterSpacing: 0,
					},
					Spacing: models.SpacingConfig{
						Margins:     "narrow",
						SectionGap:  20,
						ElementGap:  10,
						PagePadding: 32,
					},
					Decoration: models.DecorationConfig{
						Dividers:    "none",
						BulletStyle: "dash",
						UseIcons:    false,
						BorderStyle: "none",
					},
				},
				Section: models.SectionConfig{
					Personal:       models.SectionSettings{Enabled: true, Label: "Contact", LayoutType: "list"},
					Summary:        models.SectionSettings{Enabled: true, Label: "Summary", LayoutType: "list"},
					Experience:     models.SectionSettings{Enabled: true, Label: "Experience", LayoutType: "list"},
					Education:      models.SectionSettings{Enabled: true, Label: "Education", LayoutType: "list"},
					Skills:         models.SectionSettings{Enabled: true, Label: "Skills", LayoutType: "list"},
					Projects:       models.SectionSettings{Enabled: false, Label: "Projects", LayoutType: "list"},
					Certifications: models.SectionSettings{Enabled: true, Label: "Certifications", LayoutType: "list"},
				},
			},
		},
		{
			Name:        "Simple Lines",
			Category:    models.CategoryMinimalist,
			Description: "Subtle line dividers, generous whitespace. Elegant and professional.",
			IsPremium:   false,
			ATSScore:    94,
			BestFor:     models.StringArray{"Engineer", "Analyst", "Researcher"},
			Config: &models.TemplateConfig{
				Layout: models.LayoutConfig{
					Columns:       1,
					PhotoPosition: "none",
					SectionOrder:  []string{"personal", "summary", "experience", "education", "skills", "certifications"},
					PageBreak:     "auto",
				},
				Style: models.StyleConfig{
					Colors: models.ColorConfig{
						Primary:    "#374151",
						Secondary:  "#4b5563",
						Text:       "#4b5563",
						Background: "#ffffff",
						Accent:     "#6b7280",
					},
					Typography: models.TypographyConfig{
						HeadingFont:   "Inter",
						BodyFont:      "Inter",
						BaseFontSize:  10,
						LineHeight:    1.5,
						LetterSpacing: 0,
					},
					Spacing: models.SpacingConfig{
						Margins:     "normal",
						SectionGap:  22,
						ElementGap:  12,
						PagePadding: 44,
					},
					Decoration: models.DecorationConfig{
						Dividers:    "line",
						BulletStyle: "disc",
						UseIcons:    false,
						BorderStyle: "none",
					},
				},
				Section: models.SectionConfig{
					Personal:       models.SectionSettings{Enabled: true, Label: "Contact", LayoutType: "list"},
					Summary:        models.SectionSettings{Enabled: true, Label: "Profile", LayoutType: "list"},
					Experience:     models.SectionSettings{Enabled: true, Label: "Experience", LayoutType: "list"},
					Education:      models.SectionSettings{Enabled: true, Label: "Education", LayoutType: "list"},
					Skills:         models.SectionSettings{Enabled: true, Label: "Skills", LayoutType: "list"},
					Projects:       models.SectionSettings{Enabled: false, Label: "Projects", LayoutType: "list"},
					Certifications: models.SectionSettings{Enabled: true, Label: "Certifications", LayoutType: "list"},
				},
			},
		},
		// Tech
		{
			Name:        "Developer Focus",
			Category:    models.CategoryTech,
			Description: "Skills-first layout for developers. Highlights technologies and projects.",
			IsPremium:   false,
			ATSScore:    91,
			BestFor:     models.StringArray{"Frontend Dev", "Backend Dev", "Full Stack"},
			Config: &models.TemplateConfig{
				Layout: models.LayoutConfig{
					Columns:       2,
					PhotoPosition: "none",
					SectionOrder:  []string{"personal", "summary", "skills", "experience", "projects", "education"},
					PageBreak:     "auto",
				},
				Style: models.StyleConfig{
					Colors: models.ColorConfig{
						Primary:    "#059669",
						Secondary:  "#047857",
						Text:       "#111827",
						Background: "#ffffff",
						Accent:     "#10b981",
					},
					Typography: models.TypographyConfig{
						HeadingFont:   "Inter",
						BodyFont:      "Inter",
						BaseFontSize:  10,
						LineHeight:    1.45,
						LetterSpacing: 0,
					},
					Spacing: models.SpacingConfig{
						Margins:     "normal",
						SectionGap:  12,
						ElementGap:  6,
						PagePadding: 38,
					},
					Decoration: models.DecorationConfig{
						Dividers:    "line",
						BulletStyle: "dash",
						UseIcons:    true,
						BorderStyle: "none",
					},
				},
				Section: models.SectionConfig{
					Personal:   models.SectionSettings{Enabled: true, Label: "Contact", LayoutType: "list"},
					Summary:    models.SectionSettings{Enabled: true, Label: "About", LayoutType: "list"},
					Experience: models.SectionSettings{Enabled: true, Label: "Experience", LayoutType: "list"},
					Education:  models.SectionSettings{Enabled: true, Label: "Education", LayoutType: "list"},
					Skills:     models.SectionSettings{Enabled: true, Label: "Tech Stack", LayoutType: "grid"},
					Projects:   models.SectionSettings{Enabled: true, Label: "Projects", LayoutType: "list"},
					Certifications: models.SectionSettings{Enabled: true, Label: "Certifications", LayoutType: "list"},
				},
			},
		},
		// Executive (3)
		{
			Name:        "Executive Summary",
			Category:    models.CategoryExecutive,
			Description: "Bold header, two-column layout. For senior and C-level executives.",
			IsPremium:   false,
			ATSScore:    89,
			BestFor:     models.StringArray{"VP", "Director", "C-Suite"},
			Config: &models.TemplateConfig{
				Layout: models.LayoutConfig{
					Columns:       2,
					PhotoPosition: "right",
					SectionOrder:  []string{"personal", "summary", "experience", "education", "skills", "certifications"},
					PageBreak:     "auto",
				},
				Style: models.StyleConfig{
					Colors: models.ColorConfig{
						Primary:    "#1e3a8a",
						Secondary:  "#1e40af",
						Text:       "#1e293b",
						Background: "#ffffff",
						Accent:     "#3b82f6",
					},
					Typography: models.TypographyConfig{
						HeadingFont:   "Inter",
						BodyFont:      "Inter",
						BaseFontSize:  10,
						LineHeight:    1.4,
						LetterSpacing: 0,
					},
					Spacing: models.SpacingConfig{
						Margins:     "normal",
						SectionGap:  14,
						ElementGap:  6,
						PagePadding: 36,
					},
					Decoration: models.DecorationConfig{
						Dividers:    "line",
						BulletStyle: "disc",
						UseIcons:    false,
						BorderStyle: "thin",
					},
				},
				Section: models.SectionConfig{
					Personal:       models.SectionSettings{Enabled: true, Label: "Contact", LayoutType: "list"},
					Summary:        models.SectionSettings{Enabled: true, Label: "Executive Summary", LayoutType: "list"},
					Experience:     models.SectionSettings{Enabled: true, Label: "Leadership Experience", LayoutType: "list"},
					Education:      models.SectionSettings{Enabled: true, Label: "Education", LayoutType: "list"},
					Skills:         models.SectionSettings{Enabled: true, Label: "Core Competencies", LayoutType: "list"},
					Projects:       models.SectionSettings{Enabled: false, Label: "Projects", LayoutType: "list"},
					Certifications: models.SectionSettings{Enabled: true, Label: "Certifications", LayoutType: "list"},
				},
			},
		},
		{
			Name:        "Corporate Leader",
			Category:    models.CategoryExecutive,
			Description: "Refined navy and gray. Conveys authority and professionalism.",
			IsPremium:   false,
			ATSScore:    90,
			BestFor:     models.StringArray{"Manager", "Senior Manager", "Director"},
			Config: &models.TemplateConfig{
				Layout: models.LayoutConfig{
					Columns:       1,
					PhotoPosition: "none",
					SectionOrder:  []string{"personal", "summary", "experience", "education", "skills"},
					PageBreak:     "auto",
				},
				Style: models.StyleConfig{
					Colors: models.ColorConfig{
						Primary:    "#0f172a",
						Secondary:  "#1e293b",
						Text:       "#334155",
						Background: "#f8fafc",
						Accent:     "#475569",
					},
					Typography: models.TypographyConfig{
						HeadingFont:   "Inter",
						BodyFont:      "Inter",
						BaseFontSize:  11,
						LineHeight:    1.45,
						LetterSpacing: 0,
					},
					Spacing: models.SpacingConfig{
						Margins:     "wide",
						SectionGap:  20,
						ElementGap:  10,
						PagePadding: 48,
					},
					Decoration: models.DecorationConfig{
						Dividers:    "line",
						BulletStyle: "disc",
						UseIcons:    false,
						BorderStyle: "thin",
					},
				},
				Section: models.SectionConfig{
					Personal:   models.SectionSettings{Enabled: true, Label: "Contact", LayoutType: "list"},
					Summary:    models.SectionSettings{Enabled: true, Label: "Summary", LayoutType: "list"},
					Experience: models.SectionSettings{Enabled: true, Label: "Experience", LayoutType: "list"},
					Education:  models.SectionSettings{Enabled: true, Label: "Education", LayoutType: "list"},
					Skills:     models.SectionSettings{Enabled: true, Label: "Skills", LayoutType: "list"},
					Projects:   models.SectionSettings{Enabled: false, Label: "Projects", LayoutType: "list"},
					Certifications: models.SectionSettings{Enabled: true, Label: "Certifications", LayoutType: "list"},
				},
			},
		},
		{
			Name:        "Senior Executive",
			Category:    models.CategoryExecutive,
			Description: "Distinguished layout for seasoned professionals. Clean and impactful.",
			IsPremium:   false,
			ATSScore:    92,
			BestFor:     models.StringArray{"C-Level", "Board Member", "Partner"},
			Config: &models.TemplateConfig{
				Layout: models.LayoutConfig{
					Columns:       2,
					PhotoPosition: "left",
					SectionOrder:  []string{"personal", "summary", "experience", "education", "skills", "certifications"},
					PageBreak:     "auto",
				},
				Style: models.StyleConfig{
					Colors: models.ColorConfig{
						Primary:    "#1f2937",
						Secondary:  "#374151",
						Text:       "#4b5563",
						Background: "#ffffff",
						Accent:     "#6b7280",
					},
					Typography: models.TypographyConfig{
						HeadingFont:   "Georgia",
						BodyFont:      "Inter",
						BaseFontSize:  10,
						LineHeight:    1.45,
						LetterSpacing: 0,
					},
					Spacing: models.SpacingConfig{
						Margins:     "normal",
						SectionGap:  16,
						ElementGap:  8,
						PagePadding: 40,
					},
					Decoration: models.DecorationConfig{
						Dividers:    "line",
						BulletStyle: "disc",
						UseIcons:    false,
						BorderStyle: "medium",
					},
				},
				Section: models.SectionConfig{
					Personal:       models.SectionSettings{Enabled: true, Label: "Contact", LayoutType: "list"},
					Summary:        models.SectionSettings{Enabled: true, Label: "Executive Profile", LayoutType: "list"},
					Experience:     models.SectionSettings{Enabled: true, Label: "Career History", LayoutType: "list"},
					Education:      models.SectionSettings{Enabled: true, Label: "Education", LayoutType: "list"},
					Skills:         models.SectionSettings{Enabled: true, Label: "Key Competencies", LayoutType: "list"},
					Projects:       models.SectionSettings{Enabled: false, Label: "Projects", LayoutType: "list"},
					Certifications: models.SectionSettings{Enabled: true, Label: "Certifications", LayoutType: "list"},
				},
			},
		},
		// Academic (2)
		{
			Name:        "Academic Standard",
			Category:    models.CategoryAcademic,
			Description: "Structured for researchers and academics. Emphasis on publications and education.",
			IsPremium:   false,
			ATSScore:    86,
			BestFor:     models.StringArray{"Researcher", "Professor", "PhD Candidate"},
			Config: &models.TemplateConfig{
				Layout: models.LayoutConfig{
					Columns:       1,
					PhotoPosition: "none",
					SectionOrder:  []string{"personal", "summary", "education", "experience", "skills", "certifications"},
					PageBreak:     "auto",
				},
				Style: models.StyleConfig{
					Colors: models.ColorConfig{
						Primary:    "#1d4ed8",
						Secondary:  "#2563eb",
						Text:       "#1e293b",
						Background: "#ffffff",
						Accent:     "#60a5fa",
					},
					Typography: models.TypographyConfig{
						HeadingFont:   "Georgia",
						BodyFont:      "Georgia",
						BaseFontSize:  11,
						LineHeight:    1.5,
						LetterSpacing: 0,
					},
					Spacing: models.SpacingConfig{
						Margins:     "wide",
						SectionGap:  18,
						ElementGap:  8,
						PagePadding: 44,
					},
					Decoration: models.DecorationConfig{
						Dividers:    "line",
						BulletStyle: "disc",
						UseIcons:    false,
						BorderStyle: "none",
					},
				},
				Section: models.SectionConfig{
					Personal:       models.SectionSettings{Enabled: true, Label: "Contact", LayoutType: "list"},
					Summary:        models.SectionSettings{Enabled: true, Label: "Research Interests", LayoutType: "list"},
					Experience:     models.SectionSettings{Enabled: true, Label: "Experience", LayoutType: "list"},
					Education:      models.SectionSettings{Enabled: true, Label: "Education", LayoutType: "list"},
					Skills:         models.SectionSettings{Enabled: true, Label: "Skills", LayoutType: "list"},
					Projects:       models.SectionSettings{Enabled: false, Label: "Projects", LayoutType: "list"},
					Certifications: models.SectionSettings{Enabled: true, Label: "Certifications", LayoutType: "list"},
				},
			},
		},
		{
			Name:        "Research Focus",
			Category:    models.CategoryAcademic,
			Description: "Optimized for academic and research positions. Clear section hierarchy.",
			IsPremium:   false,
			ATSScore:    88,
			BestFor:     models.StringArray{"Postdoc", "Research Scientist", "Lecturer"},
			Config: &models.TemplateConfig{
				Layout: models.LayoutConfig{
					Columns:       1,
					PhotoPosition: "none",
					SectionOrder:  []string{"personal", "summary", "education", "experience", "skills"},
					PageBreak:     "auto",
				},
				Style: models.StyleConfig{
					Colors: models.ColorConfig{
						Primary:    "#4c1d95",
						Secondary:  "#5b21b6",
						Text:       "#374151",
						Background: "#ffffff",
						Accent:     "#7c3aed",
					},
					Typography: models.TypographyConfig{
						HeadingFont:   "Inter",
						BodyFont:      "Inter",
						BaseFontSize:  10,
						LineHeight:    1.5,
						LetterSpacing: 0,
					},
					Spacing: models.SpacingConfig{
						Margins:     "normal",
						SectionGap:  16,
						ElementGap:  8,
						PagePadding: 40,
					},
					Decoration: models.DecorationConfig{
						Dividers:    "line",
						BulletStyle: "circle",
						UseIcons:    false,
						BorderStyle: "none",
					},
				},
				Section: models.SectionConfig{
					Personal:   models.SectionSettings{Enabled: true, Label: "Contact", LayoutType: "list"},
					Summary:    models.SectionSettings{Enabled: true, Label: "Research Summary", LayoutType: "list"},
					Experience: models.SectionSettings{Enabled: true, Label: "Experience", LayoutType: "list"},
					Education:  models.SectionSettings{Enabled: true, Label: "Education", LayoutType: "list"},
					Skills:     models.SectionSettings{Enabled: true, Label: "Technical Skills", LayoutType: "list"},
					Projects:   models.SectionSettings{Enabled: false, Label: "Projects", LayoutType: "list"},
					Certifications: models.SectionSettings{Enabled: true, Label: "Certifications", LayoutType: "list"},
				},
			},
		},
		// Executive Sidebar - dark sidebar with nav titles, professional two-column
		{
			Name:        "Executive Sidebar",
			Category:    models.CategoryExecutive,
			Description: "Dark sidebar with section navigation. Photo, Profil, Experience, Education. Professional and distinctive.",
			IsPremium:   false,
			ATSScore:    90,
			BestFor:     models.StringArray{"Researcher", "Academic", "Executive", "Consultant"},
			Config: &models.TemplateConfig{
				Layout: models.LayoutConfig{
					Columns:       2,
					PhotoPosition: "left",
					SectionOrder:  []string{"personal", "summary", "experience", "education"},
					PageBreak:     "auto",
					LayoutVariant: "sidebar_dark",
				},
				Style: models.StyleConfig{
					Colors: models.ColorConfig{
						Primary:    "#1e3a5f",
						Secondary:  "#374151",
						Text:       "#1f2937",
						Background: "#ffffff",
						Accent:     "#1e3a5f",
					},
					Typography: models.TypographyConfig{
						HeadingFont:   "Inter",
						BodyFont:      "Inter",
						BaseFontSize:  10,
						LineHeight:    1.5,
						LetterSpacing: 0,
					},
					Spacing: models.SpacingConfig{
						Margins:     "normal",
						SectionGap:  18,
						ElementGap:  8,
						PagePadding: 0,
					},
					Decoration: models.DecorationConfig{
						Dividers:    "line",
						BulletStyle: "disc",
						UseIcons:    true,
						BorderStyle: "none",
					},
				},
				Section: models.SectionConfig{
					Personal:   models.SectionSettings{Enabled: true, Label: "Contact", LayoutType: "list"},
					Summary:    models.SectionSettings{Enabled: true, Label: "Profil", LayoutType: "list"},
					Experience: models.SectionSettings{Enabled: true, Label: "Experience", LayoutType: "list"},
					Education:  models.SectionSettings{Enabled: true, Label: "Education", LayoutType: "list"},
					Skills:     models.SectionSettings{Enabled: true, Label: "Skills", LayoutType: "list"},
					Projects:   models.SectionSettings{Enabled: false, Label: "Projects", LayoutType: "list"},
					Certifications: models.SectionSettings{Enabled: true, Label: "Certifications", LayoutType: "list"},
				},
			},
		},
	}

	for _, template := range templates {
		// Check if template already exists by name
		var existing models.Template
		if err := r.db.Where("name = ?", template.Name).First(&existing).Error; err == gorm.ErrRecordNotFound {
			if err := r.db.Create(&template).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
