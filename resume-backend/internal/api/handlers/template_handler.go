package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/resume-builder/backend/internal/repository"
)

type TemplateHandler struct {
	templateRepo *repository.TemplateRepository
}

func NewTemplateHandler(templateRepo *repository.TemplateRepository) *TemplateHandler {
	return &TemplateHandler{
		templateRepo: templateRepo,
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
