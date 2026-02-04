package repository

import (
	"github.com/resume-builder/backend/internal/config"
	"github.com/resume-builder/backend/internal/models"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(cfg *config.Config) (*gorm.DB, error) {
	logLevel := logger.Silent
	if cfg.IsDevelopment() {
		logLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}

	// Auto-migrate models
	if err := db.AutoMigrate(
		&models.User{},
		&models.Resume{},
		&models.ResumeVersion{},
		&models.Template{},
		&models.Job{},
		&models.Export{},
		&models.AIRequest{},
		&models.AuditLog{},
	); err != nil {
		return nil, err
	}

	log.Info().Msg("Database connection established and migrations applied")
	return db, nil
}
