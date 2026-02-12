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
	return `<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<style>
		@import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap');
		
		* {
			margin: 0;
			padding: 0;
			box-sizing: border-box;
		}
		
		body {
			font-family: '{{.Config.Style.Typography.BodyFont}}', 'Inter', sans-serif;
			font-size: {{.Config.Style.Typography.BaseFontSize}}pt;
			line-height: {{.Config.Style.Typography.LineHeight}};
			color: {{.Config.Style.Colors.Text}};
			background: {{.Config.Style.Colors.Background}};
		}
		
		.container {
			max-width: 8.5in;
			margin: 0 auto;
			padding: {{.Config.Style.Spacing.PagePadding}}px;
		}
		
		.header {
			text-align: center;
			margin-bottom: {{.Config.Style.Spacing.SectionGap}}px;
			padding-bottom: {{.Config.Style.Spacing.SectionGap}}px;
			border-bottom: 2px solid {{.Config.Style.Colors.Primary}};
		}
		
		.name {
			font-size: 24pt;
			font-weight: 700;
			color: {{.Config.Style.Colors.Text}};
			margin-bottom: 4px;
		}
		
		.title {
			font-size: 14pt;
			color: {{.Config.Style.Colors.Primary}};
			margin-bottom: 8px;
		}
		
		.contact-info {
			display: flex;
			justify-content: center;
			flex-wrap: wrap;
			gap: 16px;
			font-size: 10pt;
			color: #6b7280;
		}
		
		.section {
			margin-bottom: {{.Config.Style.Spacing.SectionGap}}px;
		}
		
		.section-title {
			font-size: 12pt;
			font-weight: 700;
			text-transform: uppercase;
			letter-spacing: 0.05em;
			color: {{.Config.Style.Colors.Text}};
			margin-bottom: 8px;
			padding-bottom: 4px;
			border-bottom: 1px solid #e5e7eb;
		}
		
		.summary-text {
			white-space: pre-line;
		}
		
		.entry {
			margin-bottom: {{.Config.Style.Spacing.ElementGap}}px;
		}
		
		.entry-header {
			display: flex;
			justify-content: space-between;
			align-items: flex-start;
		}
		
		.entry-title {
			font-weight: 600;
			color: {{.Config.Style.Colors.Text}};
		}
		
		.entry-subtitle {
			color: #6b7280;
		}
		
		.entry-date {
			font-size: 10pt;
			color: #6b7280;
			white-space: nowrap;
		}
		
		.entry-description {
			margin-top: 4px;
			color: #374151;
			white-space: pre-line;
		}
		
		.skills-list {
			display: flex;
			flex-wrap: wrap;
			gap: 8px;
		}
		
		.skill-tag {
			background: #f3f4f6;
			padding: 4px 12px;
			border-radius: 4px;
			font-size: 10pt;
		}
		
		{{if .Watermark}}
		.watermark {
			position: fixed;
			bottom: 20px;
			right: 20px;
			font-size: 8pt;
			color: #9ca3af;
		}
		{{end}}
	</style>
</head>
<body>
	<div class="container">
		<!-- Header -->
		<div class="header">
			<h1 class="name">{{.Data.PersonalInfo.FirstName}} {{.Data.PersonalInfo.LastName}}</h1>
			{{if .Data.PersonalInfo.Title}}
			<p class="title">{{.Data.PersonalInfo.Title}}</p>
			{{end}}
			<div class="contact-info">
				{{if .Data.PersonalInfo.Email}}<span>{{.Data.PersonalInfo.Email}}</span>{{end}}
				{{if .Data.PersonalInfo.Phone}}<span>{{.Data.PersonalInfo.Phone}}</span>{{end}}
				{{if .Data.PersonalInfo.Location}}<span>{{.Data.PersonalInfo.Location}}</span>{{end}}
				{{if .Data.PersonalInfo.LinkedIn}}<span>{{.Data.PersonalInfo.LinkedIn}}</span>{{end}}
			</div>
		</div>
		
		<!-- Summary -->
		{{if .Data.Summary}}
		<div class="section">
			<h2 class="section-title">Professional Summary</h2>
			<p class="summary-text">{{.Data.Summary}}</p>
		</div>
		{{end}}
		
		<!-- Experience -->
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
					<div class="entry-date">
						{{.StartDate}} - {{if .IsCurrent}}Present{{else}}{{.EndDate}}{{end}}
					</div>
				</div>
				{{if .Description}}
				<p class="entry-description">{{.Description}}</p>
				{{end}}
			</div>
			{{end}}
		</div>
		{{end}}
		
		<!-- Education -->
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
					<div class="entry-date">
						{{.StartDate}} - {{.EndDate}}
					</div>
				</div>
			</div>
			{{end}}
		</div>
		{{end}}
		
		<!-- Skills -->
		{{if .Data.Skills.Technical}}
		<div class="section">
			<h2 class="section-title">Skills</h2>
			<div class="skills-list">
				{{range .Data.Skills.Technical}}
				<span class="skill-tag">{{.Name}}</span>
				{{end}}
			</div>
		</div>
		{{end}}
	</div>
	
	{{if .Watermark}}
	<div class="watermark">Created with ResumeBuilder</div>
	{{end}}
</body>
</html>`
}
