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
