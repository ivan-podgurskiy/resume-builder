package repository

import (
	"github.com/google/uuid"
	"github.com/resume-builder/backend/internal/models"
	"gorm.io/gorm"
)

type ResumeRepository struct {
	db *gorm.DB
}

func NewResumeRepository(db *gorm.DB) *ResumeRepository {
	return &ResumeRepository{db: db}
}

func (r *ResumeRepository) Create(resume *models.Resume) error {
	return r.db.Create(resume).Error
}

func (r *ResumeRepository) FindByID(id uuid.UUID) (*models.Resume, error) {
	var resume models.Resume
	err := r.db.Preload("Template").First(&resume, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &resume, nil
}

func (r *ResumeRepository) FindByIDAndUser(id, userID uuid.UUID) (*models.Resume, error) {
	var resume models.Resume
	err := r.db.Preload("Template").First(&resume, "id = ? AND user_id = ?", id, userID).Error
	if err != nil {
		return nil, err
	}
	return &resume, nil
}

func (r *ResumeRepository) FindAllByUser(userID uuid.UUID, limit, offset int) ([]models.Resume, int64, error) {
	var resumes []models.Resume
	var total int64

	// Count total
	if err := r.db.Model(&models.Resume{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := r.db.Where("user_id = ?", userID).
		Order("updated_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&resumes).Error

	return resumes, total, err
}

func (r *ResumeRepository) Update(resume *models.Resume) error {
	return r.db.Save(resume).Error
}

func (r *ResumeRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Resume{}, "id = ?", id).Error
}

func (r *ResumeRepository) FindPublicBySlug(slug string) (*models.Resume, error) {
	var resume models.Resume
	err := r.db.Preload("Template").First(&resume, "public_slug = ? AND is_public = true", slug).Error
	if err != nil {
		return nil, err
	}
	return &resume, nil
}

func (r *ResumeRepository) FindPublicByID(id uuid.UUID) (*models.Resume, error) {
	var resume models.Resume
	err := r.db.Preload("Template").First(&resume, "id = ? AND is_public = true", id).Error
	if err != nil {
		return nil, err
	}
	return &resume, nil
}

func (r *ResumeRepository) SetPublic(id uuid.UUID, isPublic bool, slug *string) error {
	updates := map[string]interface{}{
		"is_public":   isPublic,
		"public_slug": slug,
	}
	return r.db.Model(&models.Resume{}).Where("id = ?", id).Updates(updates).Error
}

func (r *ResumeRepository) SlugExists(slug string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Resume{}).Where("public_slug = ?", slug).Count(&count).Error
	return count > 0, err
}

func (r *ResumeRepository) Duplicate(original *models.Resume, newTitle string) (*models.Resume, error) {
	newResume := &models.Resume{
		UserID:      original.UserID,
		Title:       newTitle,
		IsMaster:    false,
		TemplateID:  original.TemplateID,
		Data:        original.Data,
		StyleConfig: original.StyleConfig,
		IsPublic:    false,
	}
	if err := r.db.Create(newResume).Error; err != nil {
		return nil, err
	}
	return newResume, nil
}

// Version management
func (r *ResumeRepository) CreateVersion(version *models.ResumeVersion) error {
	return r.db.Create(version).Error
}

func (r *ResumeRepository) GetVersions(resumeID uuid.UUID) ([]models.ResumeVersion, error) {
	var versions []models.ResumeVersion
	err := r.db.Where("resume_id = ?", resumeID).
		Order("version_number DESC").
		Find(&versions).Error
	return versions, err
}

func (r *ResumeRepository) GetLatestVersionNumber(resumeID uuid.UUID) (int, error) {
	var version models.ResumeVersion
	err := r.db.Where("resume_id = ?", resumeID).
		Order("version_number DESC").
		First(&version).Error
	if err == gorm.ErrRecordNotFound {
		return 0, nil
	}
	return version.VersionNumber, err
}

func (r *ResumeRepository) GetVersion(resumeID uuid.UUID, versionID uuid.UUID) (*models.ResumeVersion, error) {
	var version models.ResumeVersion
	err := r.db.First(&version, "id = ? AND resume_id = ?", versionID, resumeID).Error
	if err != nil {
		return nil, err
	}
	return &version, nil
}
