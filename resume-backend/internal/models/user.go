package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionTier string

const (
	TierFree       SubscriptionTier = "free"
	TierPro        SubscriptionTier = "pro"
	TierEnterprise SubscriptionTier = "enterprise"
)

type SubscriptionStatus string

const (
	StatusActive   SubscriptionStatus = "active"
	StatusCanceled SubscriptionStatus = "canceled"
	StatusPastDue  SubscriptionStatus = "past_due"
	StatusTrialing SubscriptionStatus = "trialing"
)

type User struct {
	ID                    uuid.UUID          `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email                 string             `gorm:"uniqueIndex;not null;size:255" json:"email"`
	PasswordHash          string             `gorm:"not null" json:"-"`
	Name                  string             `gorm:"size:255" json:"name"`
	SubscriptionTier      SubscriptionTier   `gorm:"type:varchar(20);default:'free'" json:"subscription_tier"`
	SubscriptionStatus    SubscriptionStatus `gorm:"type:varchar(20);default:'active'" json:"subscription_status"`
	SubscriptionExpiresAt *time.Time         `json:"subscription_expires_at,omitempty"`
	StripeCustomerID      *string            `gorm:"size:255" json:"-"`
	EmailVerified         bool               `gorm:"default:false" json:"email_verified"`
	EmailVerifyToken      *string            `gorm:"size:255" json:"-"`
	PasswordResetToken    *string            `gorm:"size:255" json:"-"`
	PasswordResetExpires  *time.Time         `json:"-"`
	LastLoginAt           *time.Time         `json:"last_login_at,omitempty"`
	CreatedAt             time.Time          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt             time.Time          `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt             gorm.DeletedAt     `gorm:"index" json:"-"`

	// Relations
	Resumes []Resume `gorm:"foreignKey:UserID" json:"resumes,omitempty"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (u *User) IsPro() bool {
	return u.SubscriptionTier == TierPro || u.SubscriptionTier == TierEnterprise
}

func (u *User) IsEnterprise() bool {
	return u.SubscriptionTier == TierEnterprise
}

// Rate limits based on tier
func (u *User) GetAIImprovementsLimit() int {
	switch u.SubscriptionTier {
	case TierPro, TierEnterprise:
		return -1 // unlimited
	default:
		return 10 // 10 per hour for free tier
	}
}

func (u *User) GetAIExtractionsLimit() int {
	switch u.SubscriptionTier {
	case TierPro, TierEnterprise:
		return -1 // unlimited
	default:
		return 5 // 5 per day for free tier
	}
}

func (u *User) GetMaxResumes() int {
	switch u.SubscriptionTier {
	case TierPro, TierEnterprise:
		return -1 // unlimited
	default:
		return 5 // 1 for free tier
	}
}

func (u *User) GetPDFExportsLimit() int {
	switch u.SubscriptionTier {
	case TierPro, TierEnterprise:
		return -1 // unlimited
	default:
		return 3 // 3 per day for free tier
	}
}
