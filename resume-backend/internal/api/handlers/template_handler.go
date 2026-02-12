package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/resume-builder/backend/internal/repository"
	"github.com/resume-builder/backend/internal/service/pdf"
)

type TemplateHandler struct {
	templateRepo *repository.TemplateRepository
	pdfService   *pdf.PDFService
}

func NewTemplateHandler(templateRepo *repository.TemplateRepository, pdfService *pdf.PDFService) *TemplateHandler {
	return &TemplateHandler{
		templateRepo: templateRepo,
		pdfService:   pdfService,
	}
}

// GET /api/v1/templates
func (h *TemplateHandler) List(c *fiber.Ctx) error {
	templates, err := h.templateRepo.FindAll(false)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list templates",
		})
	}

	return c.JSON(fiber.Map{
		"templates": templates,
	})
}

// GET /api/v1/templates/:id
func (h *TemplateHandler) GetByID(c *fiber.Ctx) error {
	templateID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid template ID",
		})
	}

	template, err := h.templateRepo.FindByID(templateID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Template not found",
		})
	}

	return c.JSON(template)
}

// GET /api/v1/templates/:id/preview
func (h *TemplateHandler) GetPreview(c *fiber.Ctx) error {
	templateID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid template ID",
		})
	}

	template, err := h.templateRepo.FindByID(templateID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Template not found",
		})
	}

	// Return template config for preview rendering
	return c.JSON(fiber.Map{
		"id":     template.ID,
		"name":   template.Name,
		"config": template.Config,
	})
}

// GET /api/v1/templates/:id/preview-image
func (h *TemplateHandler) GetPreviewImage(c *fiber.Ctx) error {
	templateID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid template ID",
		})
	}

	template, err := h.templateRepo.FindByID(templateID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Template not found",
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	pngBuf, err := h.pdfService.GeneratePreviewImage(ctx, template)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate preview",
		})
	}

	c.Set("Content-Type", "image/png")
	c.Set("Cache-Control", "public, max-age=3600")
	return c.Send(pngBuf)
}
