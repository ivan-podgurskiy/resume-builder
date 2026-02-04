package resume

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/resume-builder/backend/internal/models"
	"github.com/resume-builder/backend/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrResumeNotFound  = errors.New("resume not found")
	ErrUnauthorized    = errors.New("unauthorized access to resume")
	ErrResumeLimitReached = errors.New("resume limit reached for your subscription tier")
	ErrSlugExists      = errors.New("slug already exists")
)

type ResumeService struct {
	resumeRepo   *repository.ResumeRepository
	userRepo     *repository.UserRepository
	templateRepo *repository.TemplateRepository
}

func NewResumeService(
	resumeRepo *repository.ResumeRepository,
	userRepo *repository.UserRepository,
	templateRepo *repository.TemplateRepository,
) *ResumeService {
	return &ResumeService{
		resumeRepo:   resumeRepo,
		userRepo:     userRepo,
		templateRepo: templateRepo,
	}
}

type CreateResumeInput struct {
	Title      string            `json:"title" validate:"required,min=1,max=255"`
	TemplateID *uuid.UUID        `json:"template_id"`
	Data       *models.ResumeData `json:"data"`
	IsMaster   bool              `json:"is_master"`
}

type UpdateResumeInput struct {
	Title       *string            `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	TemplateID  *uuid.UUID         `json:"template_id,omitempty"`
	Data        *models.ResumeData `json:"data,omitempty"`
	StyleConfig models.JSONB       `json:"style_config,omitempty"`
}

type ListResumesInput struct {
	Page     int `json:"page" validate:"min=1"`
	PageSize int `json:"page_size" validate:"min=1,max=100"`
}

type ListResumesOutput struct {
	Resumes []models.Resume `json:"resumes"`
	Total   int64           `json:"total"`
	Page    int             `json:"page"`
	Pages   int             `json:"pages"`
}

func (s *ResumeService) Create(userID uuid.UUID, input CreateResumeInput) (*models.Resume, error) {
	// Check user resume limit
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	count, err := s.userRepo.CountResumes(userID)
	if err != nil {
		return nil, err
	}

	maxResumes := user.GetMaxResumes()
	if maxResumes > 0 && count >= int64(maxResumes) {
		return nil, ErrResumeLimitReached
	}

	// Validate template if provided
	if input.TemplateID != nil {
		_, err := s.templateRepo.FindByID(*input.TemplateID)
		if err != nil {
			return nil, fmt.Errorf("invalid template ID: %w", err)
		}
	}

	// Create resume
	resume := &models.Resume{
		UserID:     userID,
		Title:      input.Title,
		IsMaster:   input.IsMaster,
		TemplateID: input.TemplateID,
		Data:       input.Data,
	}

	// Set default data if not provided
	if resume.Data == nil {
		resume.Data = &models.ResumeData{
			PersonalInfo: models.PersonalInfo{},
			Experience:   []models.Experience{},
			Education:    []models.Education{},
			Skills:       models.Skills{},
		}
	}

	if err := s.resumeRepo.Create(resume); err != nil {
		return nil, err
	}

	return resume, nil
}

func (s *ResumeService) GetByID(userID, resumeID uuid.UUID) (*models.Resume, error) {
	resume, err := s.resumeRepo.FindByIDAndUser(resumeID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrResumeNotFound
		}
		return nil, err
	}
	return resume, nil
}

func (s *ResumeService) List(userID uuid.UUID, input ListResumesInput) (*ListResumesOutput, error) {
	if input.Page < 1 {
		input.Page = 1
	}
	if input.PageSize < 1 {
		input.PageSize = 10
	}
	if input.PageSize > 100 {
		input.PageSize = 100
	}

	offset := (input.Page - 1) * input.PageSize
	resumes, total, err := s.resumeRepo.FindAllByUser(userID, input.PageSize, offset)
	if err != nil {
		return nil, err
	}

	pages := int(total) / input.PageSize
	if int(total)%input.PageSize > 0 {
		pages++
	}

	return &ListResumesOutput{
		Resumes: resumes,
		Total:   total,
		Page:    input.Page,
		Pages:   pages,
	}, nil
}

func (s *ResumeService) Update(userID, resumeID uuid.UUID, input UpdateResumeInput) (*models.Resume, error) {
	resume, err := s.GetByID(userID, resumeID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if input.Title != nil {
		resume.Title = *input.Title
	}
	if input.TemplateID != nil {
		// Validate template
		_, err := s.templateRepo.FindByID(*input.TemplateID)
		if err != nil {
			return nil, fmt.Errorf("invalid template ID: %w", err)
		}
		resume.TemplateID = input.TemplateID
	}
	if input.Data != nil {
		resume.Data = input.Data
	}
	if input.StyleConfig != nil {
		resume.StyleConfig = input.StyleConfig
	}

	if err := s.resumeRepo.Update(resume); err != nil {
		return nil, err
	}

	return resume, nil
}

func (s *ResumeService) Delete(userID, resumeID uuid.UUID) error {
	// Verify ownership
	_, err := s.GetByID(userID, resumeID)
	if err != nil {
		return err
	}

	return s.resumeRepo.Delete(resumeID)
}

func (s *ResumeService) Duplicate(userID, resumeID uuid.UUID, newTitle string) (*models.Resume, error) {
	original, err := s.GetByID(userID, resumeID)
	if err != nil {
		return nil, err
	}

	// Check limit
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	count, err := s.userRepo.CountResumes(userID)
	if err != nil {
		return nil, err
	}

	maxResumes := user.GetMaxResumes()
	if maxResumes > 0 && count >= int64(maxResumes) {
		return nil, ErrResumeLimitReached
	}

	if newTitle == "" {
		newTitle = original.Title + " (Copy)"
	}

	return s.resumeRepo.Duplicate(original, newTitle)
}

func (s *ResumeService) SetVisibility(userID, resumeID uuid.UUID, isPublic bool, customSlug *string) (*models.Resume, error) {
	resume, err := s.GetByID(userID, resumeID)
	if err != nil {
		return nil, err
	}

	var slug *string
	if isPublic {
		if customSlug != nil && *customSlug != "" {
			// Check if custom slug exists
			exists, err := s.resumeRepo.SlugExists(*customSlug)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, ErrSlugExists
			}
			slug = customSlug
		} else {
			// Generate random slug
			randomSlug, err := generateSlug(8)
			if err != nil {
				return nil, err
			}
			slug = &randomSlug
		}
	}

	if err := s.resumeRepo.SetPublic(resumeID, isPublic, slug); err != nil {
		return nil, err
	}

	resume.IsPublic = isPublic
	resume.PublicSlug = slug
	return resume, nil
}

func (s *ResumeService) GetPublicResume(slugOrID string) (*models.Resume, error) {
	// Try as slug first
	resume, err := s.resumeRepo.FindPublicBySlug(slugOrID)
	if err == nil {
		return resume, nil
	}

	// Try as UUID
	id, err := uuid.Parse(slugOrID)
	if err != nil {
		return nil, ErrResumeNotFound
	}

	resume, err = s.resumeRepo.FindPublicByID(id)
	if err != nil {
		return nil, ErrResumeNotFound
	}

	return resume, nil
}

// Version management
func (s *ResumeService) CreateVersion(userID, resumeID uuid.UUID, description string) (*models.ResumeVersion, error) {
	resume, err := s.GetByID(userID, resumeID)
	if err != nil {
		return nil, err
	}

	versionNum, err := s.resumeRepo.GetLatestVersionNumber(resumeID)
	if err != nil {
		return nil, err
	}

	version := &models.ResumeVersion{
		ResumeID:          resumeID,
		VersionNumber:     versionNum + 1,
		DataSnapshot:      resume.Data,
		ChangeDescription: description,
		CreatedBy:         userID,
	}

	if err := s.resumeRepo.CreateVersion(version); err != nil {
		return nil, err
	}

	return version, nil
}

func (s *ResumeService) GetVersions(userID, resumeID uuid.UUID) ([]models.ResumeVersion, error) {
	// Verify ownership
	_, err := s.GetByID(userID, resumeID)
	if err != nil {
		return nil, err
	}

	return s.resumeRepo.GetVersions(resumeID)
}

func (s *ResumeService) RestoreVersion(userID, resumeID, versionID uuid.UUID) (*models.Resume, error) {
	resume, err := s.GetByID(userID, resumeID)
	if err != nil {
		return nil, err
	}

	version, err := s.resumeRepo.GetVersion(resumeID, versionID)
	if err != nil {
		return nil, err
	}

	// Create a version of current state before restoring
	_, _ = s.CreateVersion(userID, resumeID, "Auto-save before restore")

	// Restore data
	resume.Data = version.DataSnapshot
	if err := s.resumeRepo.Update(resume); err != nil {
		return nil, err
	}

	return resume, nil
}

func generateSlug(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:length], nil
}
