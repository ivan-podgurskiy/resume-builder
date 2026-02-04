package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/resume-builder/backend/internal/api/middleware"
	"github.com/resume-builder/backend/internal/models"
	"github.com/resume-builder/backend/internal/service/ai"
)

type AIHandler struct {
	aiService *ai.AIService
	validate  *validator.Validate
}

func NewAIHandler(aiService *ai.AIService) *AIHandler {
	return &AIHandler{
		aiService: aiService,
		validate:  validator.New(),
	}
}

// POST /api/v1/ai/extract
func (h *AIHandler) ExtractData(c *fiber.Ctx) error {
	var input struct {
		Text string `json:"text" validate:"required"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Text is required",
		})
	}

	result, err := h.aiService.ExtractResumeData(c.Context(), input.Text)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to extract data",
		})
	}

	return c.JSON(result)
}

// POST /api/v1/ai/improve
func (h *AIHandler) ImproveText(c *fiber.Ctx) error {
	var input struct {
		Text    string `json:"text" validate:"required"`
		Context string `json:"context"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Text is required",
		})
	}

	result, err := h.aiService.ImproveText(c.Context(), input.Text, input.Context)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to improve text",
		})
	}

	return c.JSON(result)
}

// POST /api/v1/ai/generate-summary
func (h *AIHandler) GenerateSummary(c *fiber.Ctx) error {
	_ = middleware.GetUserID(c)

	var input struct {
		Data *models.ResumeData `json:"data" validate:"required"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if input.Data == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Resume data is required",
		})
	}

	summary, err := h.aiService.GenerateSummary(c.Context(), input.Data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate summary",
		})
	}

	return c.JSON(fiber.Map{
		"summary": summary,
	})
}

// POST /api/v1/ai/analyze-job
func (h *AIHandler) AnalyzeJob(c *fiber.Ctx) error {
	var input struct {
		ResumeData     *models.ResumeData `json:"resume_data" validate:"required"`
		JobDescription string             `json:"job_description" validate:"required"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Resume data and job description are required",
		})
	}

	result, err := h.aiService.AnalyzeJobMatch(c.Context(), input.ResumeData, input.JobDescription)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to analyze job match",
		})
	}

	return c.JSON(result)
}
