package pdf

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"os"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/resume-builder/backend/internal/models"
	"github.com/rs/zerolog/log"
)

// SampleResumeData returns sample data for template previews
func SampleResumeData() *models.ResumeData {
	return &models.ResumeData{
		PersonalInfo: models.PersonalInfo{
			FirstName: "Alex",
			LastName:  "Johnson",
			Title:     "Senior Software Engineer",
			Email:     "alex.johnson@email.com",
			Phone:     "(555) 123-4567",
			Location:  "San Francisco, CA",
			LinkedIn:  "linkedin.com/in/alexjohnson",
		},
		Summary: "Experienced software engineer with 8+ years building scalable applications. Passionate about clean code, system design, and mentoring junior developers.",
		Experience: []models.Experience{
			{
				ID:          "1",
				Company:     "Tech Corp",
				Position:    "Senior Software Engineer",
				Location:    "San Francisco, CA",
				StartDate:   "2020",
				EndDate:     "",
				IsCurrent:   true,
				Description: "Led development of microservices architecture. Mentored team of 5 engineers.",
				Order:       0,
			},
			{
				ID:          "2",
				Company:     "StartupXYZ",
				Position:    "Software Engineer",
				Location:    "Remote",
				StartDate:   "2017",
				EndDate:     "2020",
				IsCurrent:   false,
				Description: "Built REST APIs and real-time features. Improved deployment pipeline.",
				Order:       1,
			},
		},
		Education: []models.Education{
			{
				ID:          "1",
				Institution:  "State University",
				Degree:      "B.S.",
				FieldOfStudy: "Computer Science",
				Location:    "Boston, MA",
				StartDate:   "2013",
				EndDate:     "2017",
				Order:       0,
			},
		},
		Skills: models.Skills{
			Technical: []models.Skill{
				{Name: "JavaScript", Level: "Expert"},
				{Name: "TypeScript", Level: "Advanced"},
				{Name: "React", Level: "Expert"},
				{Name: "Node.js", Level: "Advanced"},
				{Name: "PostgreSQL", Level: "Intermediate"},
				{Name: "AWS", Level: "Intermediate"},
			},
		},
	}
}

type PDFService struct{}

func NewPDFService() *PDFService {
	return &PDFService{}
}

type PDFOptions struct {
	Quality    string // print, screen, draft
	Watermark  bool
	PageSize   string // letter, a4
	Orientation string // portrait, landscape
}

func DefaultPDFOptions() PDFOptions {
	return PDFOptions{
		Quality:    "print",
		Watermark:  false,
		PageSize:   "letter",
		Orientation: "portrait",
	}
}

// GeneratePDF creates a PDF from resume data using headless Chrome
func (s *PDFService) GeneratePDF(ctx context.Context, resume *models.Resume, template *models.Template, options PDFOptions) ([]byte, error) {
	// Generate HTML from resume data
	html, err := s.generateHTML(resume, template, options)
	if err != nil {
		return nil, fmt.Errorf("failed to generate HTML: %w", err)
	}

	// Create Chrome context with options
	opts := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
	}

	// Use custom Chrome path if specified
	if chromePath := os.Getenv("CHROME_PATH"); chromePath != "" {
		opts = append(opts, chromedp.ExecPath(chromePath))
	}

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	chromeCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	var pdfBuf []byte
	
	// Convert HTML to PDF
	if err := chromedp.Run(chromeCtx,
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}
			return page.SetDocumentContent(frameTree.Frame.ID, html).Do(ctx)
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdfBuf, _, err = page.PrintToPDF().
				WithPrintBackground(true).
				WithPaperWidth(8.5).
				WithPaperHeight(11).
				WithMarginTop(0.4).
				WithMarginBottom(0.4).
				WithMarginLeft(0.4).
				WithMarginRight(0.4).
				Do(ctx)
			return err
		}),
	); err != nil {
		log.Error().Err(err).Msg("Failed to generate PDF with chromedp")
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return pdfBuf, nil
}

// GeneratePreviewImage creates a PNG preview of a template with sample data
func (s *PDFService) GeneratePreviewImage(ctx context.Context, tmpl *models.Template) ([]byte, error) {
	sampleResume := &models.Resume{
		Data: SampleResumeData(),
	}

	html, err := s.generateHTML(sampleResume, tmpl, DefaultPDFOptions())
	if err != nil {
		return nil, fmt.Errorf("failed to generate HTML: %w", err)
	}

	opts := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
	}

	if chromePath := os.Getenv("CHROME_PATH"); chromePath != "" {
		opts = append(opts, chromedp.ExecPath(chromePath))
	}

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	chromeCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Viewport: 340x440 - crisp thumbnail for template picker (2x scale for Retina)
	var pngBuf []byte
	if err := chromedp.Run(chromeCtx,
		chromedp.EmulateViewport(340, 440),
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}
			return page.SetDocumentContent(frameTree.Frame.ID, html).Do(ctx)
		}),
		chromedp.Sleep(100),
		chromedp.CaptureScreenshot(&pngBuf),
	); err != nil {
		log.Error().Err(err).Msg("Failed to generate preview with chromedp")
		return nil, fmt.Errorf("failed to generate preview: %w", err)
	}

	return pngBuf, nil
}

func (s *PDFService) generateHTML(resume *models.Resume, tmpl *models.Template, options PDFOptions) (string, error) {
	data := resume.Data
	if data == nil {
		return "", fmt.Errorf("resume data is nil")
	}

	// Get template config or use defaults
	var config *models.TemplateConfig
	if tmpl != nil && tmpl.Config != nil {
		config = tmpl.Config
	} else {
		config = defaultTemplateConfig()
	}

	// Build HTML
	htmlTemplate := buildResumeTemplate(config)
	
	t, err := template.New("resume").Parse(htmlTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	templateData := map[string]interface{}{
		"Data":      data,
		"Config":    config,
		"Watermark": options.Watermark,
	}

	if err := t.Execute(&buf, templateData); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

func defaultTemplateConfig() *models.TemplateConfig {
	return &models.TemplateConfig{
		Layout: models.LayoutConfig{
			Columns:       1,
			PhotoPosition: "none",
			SectionOrder:  []string{"personal", "summary", "experience", "education", "skills"},
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
				HeadingFont:  "Inter",
				BodyFont:     "Inter",
				BaseFontSize: 11,
				LineHeight:   1.5,
			},
			Spacing: models.SpacingConfig{
				Margins:     "normal",
				SectionGap:  16,
				ElementGap:  8,
				PagePadding: 40,
			},
		},
	}
}

func buildResumeTemplate(config *models.TemplateConfig) string {
	sidebarDark := config.Layout.LayoutVariant == "sidebar_dark"
	twoCol := config.Layout.Columns == 2 && !sidebarDark
	layoutCss := ""
	layoutBody := ""
	if sidebarDark {
		layoutCss = `
		.resume-sidebar-dark { display: grid; grid-template-columns: 28% 1fr; min-height: 11in; }
		.sidebar-dark { background: #1f2937; color: white; padding: 24px 20px; }
		.sidebar-dark .photo-wrap { width: 80px; height: 80px; margin: 0 auto; border-radius: 50%; overflow: hidden; border: 3px solid rgba(255,255,255,0.3); background: #374151; }
		.sidebar-dark .photo-wrap img { width: 100%; height: 100%; object-fit: cover; }
		.sidebar-dark .nav-title { font-size: 10pt; font-weight: 700; text-transform: uppercase; letter-spacing: 0.1em; padding: 12px 0 12px 16px; border-top: 1px solid rgba(255,255,255,0.3); }
		.main-dark { padding: 28px 32px; position: relative; overflow: hidden; }
		.main-dark .name { font-size: 22pt; margin-bottom: 4px; }
		.main-dark .wave-deco { position: absolute; top: -20px; right: -20px; width: 120px; height: 120px; background: radial-gradient(circle at 100% 0%, #1e3a5f 0%, transparent 70%); opacity: 0.5; }
		.main-dark .section-cell { padding: 16px 0; border-top: 1px solid rgba(0,0,0,0.08); }
		.main-dark .entry-date { font-size: 9pt; }
		.main-dark .entry-description { margin-top: 6px; font-size: 9.5pt; }
		.icon-inline { display: inline-block; width: 12px; height: 12px; margin-right: 4px; vertical-align: middle; opacity: 0.7; }`
		layoutBody = `
		<div class="resume-sidebar-dark">
			<div class="sidebar-dark"><div class="photo-wrap">{{if .Data.PersonalInfo.PhotoURL}}<img src="{{.Data.PersonalInfo.PhotoURL}}" alt="Photo" />{{end}}</div></div>
			<div class="main-dark"><div class="wave-deco"></div><h1 class="name">{{.Data.PersonalInfo.FirstName}} {{.Data.PersonalInfo.LastName}}</h1>{{if .Data.PersonalInfo.Title}}<p class="title" style="color:{{.Config.Style.Colors.Primary}};">{{.Data.PersonalInfo.Title}}</p>{{end}}</div>
			<div class="sidebar-dark nav-title">Profil</div>
			<div class="main-dark section-cell">{{if .Data.Summary}}<p class="summary-text" style="margin:0;padding-left:12px;border-left:2px solid {{.Config.Style.Colors.Primary}};white-space:pre-line;">{{.Data.Summary}}</p>{{else}}<p style="color:#9ca3af;font-style:italic;">—</p>{{end}}</div>
			<div class="sidebar-dark nav-title">Experience</div>
			<div class="main-dark section-cell">{{if .Data.Experience}}{{range .Data.Experience}}<div class="entry" style="margin-bottom:12px;"><div class="entry-header" style="display:flex;justify-content:space-between;gap:12px;"><div><div class="entry-title" style="font-weight:600;">{{.Position}}</div><div class="entry-subtitle" style="font-size:9.5pt;color:#6b7280;">{{.Company}}{{if .Location}} · {{.Location}}{{end}}</div></div><div class="entry-date">{{.StartDate}} - {{if .IsCurrent}}Present{{else}}{{.EndDate}}{{end}}</div></div>{{if .Description}}<p class="entry-description">{{.Description}}</p>{{end}}</div>{{end}}{{else}}<p style="color:#9ca3af;font-style:italic;">—</p>{{end}}</div>
			<div class="sidebar-dark nav-title">Education</div>
			<div class="main-dark section-cell">{{if .Data.Education}}{{range .Data.Education}}<div class="entry" style="margin-bottom:12px;"><div class="entry-header" style="display:flex;justify-content:space-between;gap:12px;"><div><div class="entry-title" style="font-weight:600;">{{.Degree}} {{if .FieldOfStudy}}in {{.FieldOfStudy}}{{end}}</div><div class="entry-subtitle" style="font-size:9.5pt;color:#6b7280;">{{.Institution}}</div></div><div class="entry-date">{{.StartDate}} - {{.EndDate}}</div></div></div>{{end}}{{else}}<p style="color:#9ca3af;font-style:italic;">—</p>{{end}}</div>
		</div>`
	} else if twoCol {
		layoutCss = `
		.resume-grid { display: grid; grid-template-columns: 1fr 2fr; gap: 24px; align-items: start; }
		.sidebar { padding-right: 16px; border-right: 1px solid #e5e7eb; }
		.sidebar .section-title { font-size: 10pt; }
		.main-content { min-width: 0; }`
		layoutBody = `
		<div class="resume-grid">
			<div class="sidebar">
				<div class="header">
					<h1 class="name">{{.Data.PersonalInfo.FirstName}} {{.Data.PersonalInfo.LastName}}</h1>
					{{if .Data.PersonalInfo.Title}}<p class="title">{{.Data.PersonalInfo.Title}}</p>{{end}}
				</div>
				<div class="contact-info contact-sidebar">
					{{if .Data.PersonalInfo.Email}}<div>{{.Data.PersonalInfo.Email}}</div>{{end}}
					{{if .Data.PersonalInfo.Phone}}<div>{{.Data.PersonalInfo.Phone}}</div>{{end}}
					{{if .Data.PersonalInfo.Location}}<div>{{.Data.PersonalInfo.Location}}</div>{{end}}
					{{if .Data.PersonalInfo.LinkedIn}}<div>{{.Data.PersonalInfo.LinkedIn}}</div>{{end}}
				</div>
				{{if .Data.Skills.Technical}}
				<div class="section">
					<h2 class="section-title">Skills</h2>
					<div class="skills-list">
						{{range .Data.Skills.Technical}}<span class="skill-tag">{{.Name}}</span>{{end}}
					</div>
				</div>
				{{end}}
				{{if .Data.Education}}
				<div class="section">
					<h2 class="section-title">Education</h2>
					{{range .Data.Education}}
					<div class="entry">
						<div class="entry-title">{{.Degree}}</div>
						<div class="entry-subtitle">{{.Institution}}</div>
						<div class="entry-date">{{.StartDate}} - {{.EndDate}}</div>
					</div>
					{{end}}
				</div>
				{{end}}
			</div>
			<div class="main-content">
				{{if .Data.Summary}}
				<div class="section">
					<h2 class="section-title">Professional Summary</h2>
					<p class="summary-text">{{.Data.Summary}}</p>
				</div>
				{{end}}
				{{if .Data.Experience}}
				<div class="section">
					<h2 class="section-title">Experience</h2>
					{{range .Data.Experience}}
					<div class="entry">
						<div class="entry-header">
							<div>
								<div class="entry-title">{{.Position}}</div>
								<div class="entry-subtitle">{{.Company}}{{if .Location}} · {{.Location}}{{end}}</div>
							</div>
							<div class="entry-date">{{.StartDate}} - {{if .IsCurrent}}Present{{else}}{{.EndDate}}{{end}}</div>
						</div>
						{{if .Description}}<p class="entry-description">{{.Description}}</p>{{end}}
					</div>
					{{end}}
				</div>
				{{end}}
			</div>
		</div>`
	} else {
		layoutBody = `
		<div class="header">
			<h1 class="name">{{.Data.PersonalInfo.FirstName}} {{.Data.PersonalInfo.LastName}}</h1>
			{{if .Data.PersonalInfo.Title}}<p class="title">{{.Data.PersonalInfo.Title}}</p>{{end}}
			<div class="contact-info">
				{{if .Data.PersonalInfo.Email}}<span>{{.Data.PersonalInfo.Email}}</span>{{end}}
				{{if .Data.PersonalInfo.Phone}}<span>{{.Data.PersonalInfo.Phone}}</span>{{end}}
				{{if .Data.PersonalInfo.Location}}<span>{{.Data.PersonalInfo.Location}}</span>{{end}}
				{{if .Data.PersonalInfo.LinkedIn}}<span>{{.Data.PersonalInfo.LinkedIn}}</span>{{end}}
			</div>
		</div>
		{{if .Data.Summary}}
		<div class="section">
			<h2 class="section-title">Professional Summary</h2>
			<p class="summary-text">{{.Data.Summary}}</p>
		</div>
		{{end}}
		{{if .Data.Experience}}
		<div class="section">
			<h2 class="section-title">Experience</h2>
			{{range .Data.Experience}}
			<div class="entry">
				<div class="entry-header">
					<div>
						<div class="entry-title">{{.Position}}</div>
						<div class="entry-subtitle">{{.Company}}{{if .Location}} · {{.Location}}{{end}}</div>
					</div>
					<div class="entry-date">{{.StartDate}} - {{if .IsCurrent}}Present{{else}}{{.EndDate}}{{end}}</div>
				</div>
				{{if .Description}}<p class="entry-description">{{.Description}}</p>{{end}}
			</div>
			{{end}}
		</div>
		{{end}}
		{{if .Data.Education}}
		<div class="section">
			<h2 class="section-title">Education</h2>
			{{range .Data.Education}}
			<div class="entry">
				<div class="entry-header">
					<div>
						<div class="entry-title">{{.Degree}} in {{.FieldOfStudy}}</div>
						<div class="entry-subtitle">{{.Institution}}</div>
					</div>
					<div class="entry-date">{{.StartDate}} - {{.EndDate}}</div>
				</div>
			</div>
			{{end}}
		</div>
		{{end}}
		{{if .Data.Skills.Technical}}
		<div class="section">
			<h2 class="section-title">Skills</h2>
			<div class="skills-list">
				{{range .Data.Skills.Technical}}<span class="skill-tag">{{.Name}}</span>{{end}}
			</div>
		</div>
		{{end}}`
	}

	return `<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<style>
		@import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap');
		* { margin: 0; padding: 0; box-sizing: border-box; }
		body {
			font-family: '{{.Config.Style.Typography.BodyFont}}', 'Inter', sans-serif;
			font-size: {{.Config.Style.Typography.BaseFontSize}}pt;
			line-height: {{.Config.Style.Typography.LineHeight}};
			color: {{.Config.Style.Colors.Text}};
			background: {{.Config.Style.Colors.Background}};
		}
		.container { max-width: 8.5in; margin: 0 auto; padding: {{.Config.Style.Spacing.PagePadding}}px; }
		.header { margin-bottom: {{.Config.Style.Spacing.SectionGap}}px; padding-bottom: {{.Config.Style.Spacing.SectionGap}}px; border-bottom: 2px solid {{.Config.Style.Colors.Primary}}; }
		.name { font-size: 24pt; font-weight: 700; color: {{.Config.Style.Colors.Text}}; margin-bottom: 4px; }
		.title { font-size: 14pt; color: {{.Config.Style.Colors.Primary}}; margin-bottom: 8px; }
		.contact-info { display: flex; justify-content: center; flex-wrap: wrap; gap: 16px; font-size: 10pt; color: #6b7280; }
		.contact-sidebar { flex-direction: column; justify-content: flex-start; gap: 8px; }
		.section { margin-bottom: {{.Config.Style.Spacing.SectionGap}}px; }
		.section-title { font-size: 12pt; font-weight: 700; text-transform: uppercase; letter-spacing: 0.05em; color: {{.Config.Style.Colors.Text}}; margin-bottom: 8px; padding-bottom: 4px; border-bottom: 1px solid #e5e7eb; }
		.summary-text { white-space: pre-line; }
		.entry { margin-bottom: {{.Config.Style.Spacing.ElementGap}}px; }
		.entry-header { display: flex; justify-content: space-between; align-items: flex-start; }
		.entry-title { font-weight: 600; color: {{.Config.Style.Colors.Text}}; }
		.entry-subtitle { color: #6b7280; }
		.entry-date { font-size: 10pt; color: #6b7280; white-space: nowrap; }
		.entry-description { margin-top: 4px; color: #374151; white-space: pre-line; }
		.skills-list { display: flex; flex-wrap: wrap; gap: 8px; }
		.skill-tag { background: #f3f4f6; padding: 4px 12px; border-radius: 4px; font-size: 10pt; }
		` + layoutCss + `
		{{if .Watermark}}.watermark { position: fixed; bottom: 20px; right: 20px; font-size: 8pt; color: #9ca3af; }{{end}}
	</style>
</head>
<body>
	<div class="container">` + layoutBody + `
	</div>
	{{if .Watermark}}<div class="watermark">Created with ResumeBuilder</div>{{end}}
</body>
</html>`
}
