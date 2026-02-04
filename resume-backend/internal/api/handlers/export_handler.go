package handlers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/resume-builder/backend/internal/api/middleware"
	"github.com/resume-builder/backend/internal/models"
	"github.com/resume-builder/backend/internal/repository"
	"github.com/resume-builder/backend/internal/service/pdf"
)

type ExportHandler struct {
	resumeRepo   *repository.ResumeRepository
	templateRepo *repository.TemplateRepository
	pdfService   *pdf.PDFService
	validate     *validator.Validate
}

func NewExportHandler(
	resumeRepo *repository.ResumeRepository,
	templateRepo *repository.TemplateRepository,
	pdfService *pdf.PDFService,
) *ExportHandler {
	return &ExportHandler{
		resumeRepo:   resumeRepo,
		templateRepo: templateRepo,
		pdfService:   pdfService,
		validate:     validator.New(),
	}
}

// POST /api/v1/export/pdf
func (h *ExportHandler) ExportPDF(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var input struct {
		ResumeID  string `json:"resume_id" validate:"required,uuid"`
		Quality   string `json:"quality"`   // print, screen, draft
		Watermark bool   `json:"watermark"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Resume ID is required",
		})
	}

	resumeID, _ := uuid.Parse(input.ResumeID)

	// Get resume
	resume, err := h.resumeRepo.FindByIDAndUser(resumeID, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Resume not found",
		})
	}

	// Get template if set
	var template *models.Template
	if resume.TemplateID != nil {
		template, _ = h.templateRepo.FindByID(*resume.TemplateID)
	}

	// Set options
	options := pdf.DefaultPDFOptions()
	if input.Quality != "" {
		options.Quality = input.Quality
	}
	options.Watermark = input.Watermark

	// Generate PDF (synchronous for now, should be async in production)
	pdfBytes, err := h.pdfService.GeneratePDF(c.Context(), resume, template, options)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate PDF",
		})
	}

	// Set headers for PDF download
	filename := fmt.Sprintf("%s_%s.pdf", sanitizeFilename(resume.Title), time.Now().Format("20060102"))
	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Set("Content-Length", fmt.Sprintf("%d", len(pdfBytes)))

	return c.Send(pdfBytes)
}

// POST /api/v1/export/txt
func (h *ExportHandler) ExportTXT(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var input struct {
		ResumeID string `json:"resume_id" validate:"required,uuid"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	resumeID, _ := uuid.Parse(input.ResumeID)

	// Get resume
	resume, err := h.resumeRepo.FindByIDAndUser(resumeID, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Resume not found",
		})
	}

	// Generate plain text
	txt := generatePlainText(resume)

	// Set headers for TXT download
	filename := fmt.Sprintf("%s_%s.txt", sanitizeFilename(resume.Title), time.Now().Format("20060102"))
	c.Set("Content-Type", "text/plain; charset=utf-8")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	return c.SendString(txt)
}

// POST /api/v1/export/json
func (h *ExportHandler) ExportJSON(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var input struct {
		ResumeID string `json:"resume_id" validate:"required,uuid"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	resumeID, _ := uuid.Parse(input.ResumeID)

	// Get resume
	resume, err := h.resumeRepo.FindByIDAndUser(resumeID, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Resume not found",
		})
	}

	// Generate JSON
	jsonBytes, err := json.MarshalIndent(resume.Data, "", "  ")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate JSON",
		})
	}

	// Set headers for JSON download
	filename := fmt.Sprintf("%s_%s.json", sanitizeFilename(resume.Title), time.Now().Format("20060102"))
	c.Set("Content-Type", "application/json")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	return c.Send(jsonBytes)
}

func sanitizeFilename(name string) string {
	// Remove/replace characters that are problematic in filenames
	result := ""
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			result += string(r)
		} else if r == ' ' {
			result += "_"
		}
	}
	if result == "" {
		result = "resume"
	}
	return result
}

func generatePlainText(resume *models.Resume) string {
	data := resume.Data
	if data == nil {
		return ""
	}

	txt := ""

	// Header
	txt += fmt.Sprintf("%s %s\n", data.PersonalInfo.FirstName, data.PersonalInfo.LastName)
	if data.PersonalInfo.Title != "" {
		txt += fmt.Sprintf("%s\n", data.PersonalInfo.Title)
	}
	txt += "\n"

	// Contact
	txt += "CONTACT\n"
	txt += fmt.Sprintf("----------------------------------------\n")
	if data.PersonalInfo.Email != "" {
		txt += fmt.Sprintf("Email: %s\n", data.PersonalInfo.Email)
	}
	if data.PersonalInfo.Phone != "" {
		txt += fmt.Sprintf("Phone: %s\n", data.PersonalInfo.Phone)
	}
	if data.PersonalInfo.Location != "" {
		txt += fmt.Sprintf("Location: %s\n", data.PersonalInfo.Location)
	}
	if data.PersonalInfo.LinkedIn != "" {
		txt += fmt.Sprintf("LinkedIn: %s\n", data.PersonalInfo.LinkedIn)
	}
	if data.PersonalInfo.Website != "" {
		txt += fmt.Sprintf("Website: %s\n", data.PersonalInfo.Website)
	}
	txt += "\n"

	// Summary
	if data.Summary != "" {
		txt += "PROFESSIONAL SUMMARY\n"
		txt += fmt.Sprintf("----------------------------------------\n")
		txt += fmt.Sprintf("%s\n\n", data.Summary)
	}

	// Experience
	if len(data.Experience) > 0 {
		txt += "EXPERIENCE\n"
		txt += fmt.Sprintf("----------------------------------------\n")
		for _, exp := range data.Experience {
			endDate := exp.EndDate
			if exp.IsCurrent {
				endDate = "Present"
			}
			txt += fmt.Sprintf("%s\n", exp.Position)
			txt += fmt.Sprintf("%s | %s - %s\n", exp.Company, exp.StartDate, endDate)
			if exp.Description != "" {
				txt += fmt.Sprintf("%s\n", exp.Description)
			}
			txt += "\n"
		}
	}

	// Education
	if len(data.Education) > 0 {
		txt += "EDUCATION\n"
		txt += fmt.Sprintf("----------------------------------------\n")
		for _, edu := range data.Education {
			txt += fmt.Sprintf("%s in %s\n", edu.Degree, edu.FieldOfStudy)
			txt += fmt.Sprintf("%s | %s - %s\n", edu.Institution, edu.StartDate, edu.EndDate)
			txt += "\n"
		}
	}

	// Skills
	if len(data.Skills.Technical) > 0 {
		txt += "SKILLS\n"
		txt += fmt.Sprintf("----------------------------------------\n")
		skills := ""
		for i, skill := range data.Skills.Technical {
			if i > 0 {
				skills += ", "
			}
			skills += skill.Name
		}
		txt += fmt.Sprintf("%s\n", skills)
	}

	return txt
}
