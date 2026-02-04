package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/resume-builder/backend/internal/api/middleware"
	"github.com/resume-builder/backend/internal/service/resume"
)

type ResumeHandler struct {
	resumeService *resume.ResumeService
	validate      *validator.Validate
}

func NewResumeHandler(resumeService *resume.ResumeService) *ResumeHandler {
	return &ResumeHandler{
		resumeService: resumeService,
		validate:      validator.New(),
	}
}

// GET /api/v1/resumes
func (h *ResumeHandler) List(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	input := resume.ListResumesInput{
		Page:     c.QueryInt("page", 1),
		PageSize: c.QueryInt("page_size", 10),
	}

	result, err := h.resumeService.List(userID, input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list resumes",
		})
	}

	return c.JSON(result)
}

// POST /api/v1/resumes
func (h *ResumeHandler) Create(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var input resume.CreateResumeInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	result, err := h.resumeService.Create(userID, input)
	if err != nil {
		if err == resume.ErrResumeLimitReached {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Resume limit reached. Please upgrade your subscription.",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create resume",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}

// GET /api/v1/resumes/:id
func (h *ResumeHandler) GetByID(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	resumeID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid resume ID",
		})
	}

	result, err := h.resumeService.GetByID(userID, resumeID)
	if err != nil {
		if err == resume.ErrResumeNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Resume not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get resume",
		})
	}

	return c.JSON(result)
}

// PUT /api/v1/resumes/:id
func (h *ResumeHandler) Update(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	resumeID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid resume ID",
		})
	}

	var input resume.UpdateResumeInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	result, err := h.resumeService.Update(userID, resumeID, input)
	if err != nil {
		if err == resume.ErrResumeNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Resume not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update resume",
		})
	}

	return c.JSON(result)
}

// PATCH /api/v1/resumes/:id
func (h *ResumeHandler) PartialUpdate(c *fiber.Ctx) error {
	return h.Update(c)
}

// DELETE /api/v1/resumes/:id
func (h *ResumeHandler) Delete(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	resumeID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid resume ID",
		})
	}

	if err := h.resumeService.Delete(userID, resumeID); err != nil {
		if err == resume.ErrResumeNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Resume not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete resume",
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

// POST /api/v1/resumes/:id/duplicate
func (h *ResumeHandler) Duplicate(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	resumeID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid resume ID",
		})
	}

	var input struct {
		Title string `json:"title"`
	}
	_ = c.BodyParser(&input)

	result, err := h.resumeService.Duplicate(userID, resumeID, input.Title)
	if err != nil {
		if err == resume.ErrResumeNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Resume not found",
			})
		}
		if err == resume.ErrResumeLimitReached {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Resume limit reached. Please upgrade your subscription.",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to duplicate resume",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}

// GET /api/v1/resumes/:id/versions
func (h *ResumeHandler) GetVersions(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	resumeID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid resume ID",
		})
	}

	versions, err := h.resumeService.GetVersions(userID, resumeID)
	if err != nil {
		if err == resume.ErrResumeNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Resume not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get versions",
		})
	}

	return c.JSON(fiber.Map{
		"versions": versions,
	})
}

// POST /api/v1/resumes/:id/versions/:versionId/restore
func (h *ResumeHandler) RestoreVersion(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	resumeID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid resume ID",
		})
	}

	versionID, err := uuid.Parse(c.Params("versionId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid version ID",
		})
	}

	result, err := h.resumeService.RestoreVersion(userID, resumeID, versionID)
	if err != nil {
		if err == resume.ErrResumeNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Resume or version not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to restore version",
		})
	}

	return c.JSON(result)
}

// PATCH /api/v1/resumes/:id/visibility
func (h *ResumeHandler) SetVisibility(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	resumeID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid resume ID",
		})
	}

	var input struct {
		IsPublic bool    `json:"is_public"`
		Slug     *string `json:"slug,omitempty"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	result, err := h.resumeService.SetVisibility(userID, resumeID, input.IsPublic, input.Slug)
	if err != nil {
		if err == resume.ErrResumeNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Resume not found",
			})
		}
		if err == resume.ErrSlugExists {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Slug already in use",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update visibility",
		})
	}

	return c.JSON(result)
}

// GET /api/v1/resumes/:id/public (public endpoint, no auth required)
func (h *ResumeHandler) GetPublicResume(c *fiber.Ctx) error {
	slugOrID := c.Params("id")

	result, err := h.resumeService.GetPublicResume(slugOrID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Resume not found",
		})
	}

	return c.JSON(result)
}
