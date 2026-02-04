package repository

import (
	"github.com/google/uuid"
	"github.com/resume-builder/backend/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.User{}, "id = ?", id).Error
}

func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (r *UserRepository) UpdateLastLogin(id uuid.UUID) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("last_login_at", gorm.Expr("NOW()")).Error
}

func (r *UserRepository) SetEmailVerified(id uuid.UUID) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"email_verified":     true,
		"email_verify_token": nil,
	}).Error
}

func (r *UserRepository) FindByVerifyToken(token string) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "email_verify_token = ?", token).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByResetToken(token string) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "password_reset_token = ? AND password_reset_expires > NOW()", token).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) CountResumes(userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.Resume{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}
