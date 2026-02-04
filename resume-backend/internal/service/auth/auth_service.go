package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/resume-builder/backend/internal/api/middleware"
	"github.com/resume-builder/backend/internal/config"
	"github.com/resume-builder/backend/internal/models"
	"github.com/resume-builder/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailExists        = errors.New("email already exists")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrEmailNotVerified   = errors.New("email not verified")
)

type AuthService struct {
	userRepo   *repository.UserRepository
	authMiddleware *middleware.AuthMiddleware
	config     *config.Config
}

func NewAuthService(userRepo *repository.UserRepository, authMiddleware *middleware.AuthMiddleware, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		authMiddleware: authMiddleware,
		config:     cfg,
	}
}

type RegisterInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required,min=2"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	User         *models.User `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int          `json:"expires_in"` // seconds
}

func (s *AuthService) Register(input RegisterInput) (*AuthResponse, error) {
	// Check if email exists
	exists, err := s.userRepo.ExistsByEmail(input.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
	if err != nil {
		return nil, err
	}

	// Generate email verification token
	verifyToken, err := generateToken(32)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Email:            input.Email,
		PasswordHash:     string(hashedPassword),
		Name:             input.Name,
		SubscriptionTier: models.TierFree,
		EmailVerifyToken: &verifyToken,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Generate tokens
	accessToken, err := s.authMiddleware.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.authMiddleware.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.config.JWTExpirationHours * 3600,
	}, nil
}

func (s *AuthService) Login(input LoginInput) (*AuthResponse, error) {
	// Find user
	user, err := s.userRepo.FindByEmail(input.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Update last login
	_ = s.userRepo.UpdateLastLogin(user.ID)

	// Generate tokens
	accessToken, err := s.authMiddleware.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.authMiddleware.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.config.JWTExpirationHours * 3600,
	}, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (*AuthResponse, error) {
	claims, err := s.authMiddleware.ValidateToken(refreshToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Generate new tokens
	accessToken, err := s.authMiddleware.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.authMiddleware.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    s.config.JWTExpirationHours * 3600,
	}, nil
}

func (s *AuthService) GetCurrentUser(userID uuid.UUID) (*models.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *AuthService) VerifyEmail(token string) error {
	user, err := s.userRepo.FindByVerifyToken(token)
	if err != nil {
		return ErrInvalidToken
	}

	return s.userRepo.SetEmailVerified(user.ID)
}

func (s *AuthService) ForgotPassword(email string) (string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		// Don't reveal if email exists
		return "", nil
	}

	// Generate reset token
	resetToken, err := generateToken(32)
	if err != nil {
		return "", err
	}

	// Set expiration (1 hour)
	expires := time.Now().Add(time.Hour)
	user.PasswordResetToken = &resetToken
	user.PasswordResetExpires = &expires

	if err := s.userRepo.Update(user); err != nil {
		return "", err
	}

	return resetToken, nil
}

func (s *AuthService) ResetPassword(token, newPassword string) error {
	user, err := s.userRepo.FindByResetToken(token)
	if err != nil {
		return ErrInvalidToken
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedPassword)
	user.PasswordResetToken = nil
	user.PasswordResetExpires = nil

	return s.userRepo.Update(user)
}

func generateToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
