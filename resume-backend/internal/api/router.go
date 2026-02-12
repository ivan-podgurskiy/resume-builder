package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/resume-builder/backend/internal/api/handlers"
	"github.com/resume-builder/backend/internal/api/middleware"
	"github.com/resume-builder/backend/internal/config"
	"github.com/resume-builder/backend/internal/repository"
	"github.com/resume-builder/backend/internal/service/ai"
	"github.com/resume-builder/backend/internal/service/auth"
	"github.com/resume-builder/backend/internal/service/fileparser"
	"github.com/resume-builder/backend/internal/service/pdf"
	"github.com/resume-builder/backend/internal/service/resume"
	"gorm.io/gorm"
)

type Router struct {
	app            *fiber.App
	config         *config.Config
	db             *gorm.DB
	authMiddleware *middleware.AuthMiddleware
}

func NewRouter(cfg *config.Config, db *gorm.DB) *Router {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	authMiddleware := middleware.NewAuthMiddleware(cfg)

	return &Router{
		app:            app,
		config:         cfg,
		db:             db,
		authMiddleware: authMiddleware,
	}
}

func (r *Router) Setup() *fiber.App {
	// Global middleware
	r.app.Use(recover.New())
	r.app.Use(middleware.NewLoggerMiddleware())
	r.app.Use(middleware.NewCORSMiddleware(r.config))

	// Rate limiter
	rateLimiter := middleware.NewRateLimiter(r.config)
	r.app.Use(rateLimiter.Middleware())

	// Health check
	r.app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	// Initialize repositories
	userRepo := repository.NewUserRepository(r.db)
	resumeRepo := repository.NewResumeRepository(r.db)
	templateRepo := repository.NewTemplateRepository(r.db)

	// Seed default templates
	_ = templateRepo.SeedDefaultTemplates()

	// Initialize services
	authService := auth.NewAuthService(userRepo, r.authMiddleware, r.config)
	resumeService := resume.NewResumeService(resumeRepo, userRepo, templateRepo)
	aiService := ai.NewAIService(r.config)
	pdfService := pdf.NewPDFService()
	fileParserService := fileparser.NewFileParserService()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	resumeHandler := handlers.NewResumeHandler(resumeService)
	templateHandler := handlers.NewTemplateHandler(templateRepo, pdfService)
	aiHandler := handlers.NewAIHandler(aiService, fileParserService)
	exportHandler := handlers.NewExportHandler(resumeRepo, templateRepo, pdfService)

	// API routes
	api := r.app.Group("/api/v1")

	// Auth routes (public)
	authRoutes := api.Group("/auth")
	authRoutes.Post("/register", authHandler.Register)
	authRoutes.Post("/login", authHandler.Login)
	authRoutes.Post("/logout", authHandler.Logout)
	authRoutes.Post("/refresh", authHandler.Refresh)
	authRoutes.Post("/forgot-password", authHandler.ForgotPassword)
	authRoutes.Post("/reset-password", authHandler.ResetPassword)
	authRoutes.Post("/verify-email", authHandler.VerifyEmail)

	// Auth routes (protected)
	authRoutes.Get("/me", r.authMiddleware.Protected(), authHandler.GetCurrentUser)

	// Resume routes (protected)
	resumeRoutes := api.Group("/resumes", r.authMiddleware.Protected())
	resumeRoutes.Get("/", resumeHandler.List)
	resumeRoutes.Post("/", resumeHandler.Create)
	resumeRoutes.Get("/:id", resumeHandler.GetByID)
	resumeRoutes.Put("/:id", resumeHandler.Update)
	resumeRoutes.Patch("/:id", resumeHandler.PartialUpdate)
	resumeRoutes.Delete("/:id", resumeHandler.Delete)
	resumeRoutes.Post("/:id/duplicate", resumeHandler.Duplicate)
	resumeRoutes.Get("/:id/versions", resumeHandler.GetVersions)
	resumeRoutes.Post("/:id/versions/:versionId/restore", resumeHandler.RestoreVersion)
	resumeRoutes.Patch("/:id/visibility", resumeHandler.SetVisibility)

	// Public resume route (no auth)
	api.Get("/share/:id", resumeHandler.GetPublicResume)

	// Template routes (public for listing, some may be protected)
	templateRoutes := api.Group("/templates")
	templateRoutes.Get("/", templateHandler.List)
	templateRoutes.Get("/:id", templateHandler.GetByID)
	templateRoutes.Get("/:id/preview", templateHandler.GetPreview)
	templateRoutes.Get("/:id/preview-image", templateHandler.GetPreviewImage)

	// AI routes (protected)
	aiRoutes := api.Group("/ai", r.authMiddleware.Protected())
	aiRoutes.Post("/extract", aiHandler.ExtractData)
	aiRoutes.Post("/extract-file", aiHandler.ExtractFromFile)
	aiRoutes.Get("/supported-formats", aiHandler.GetSupportedFormats)
	aiRoutes.Post("/improve", aiHandler.ImproveText)
	aiRoutes.Post("/generate-summary", aiHandler.GenerateSummary)
	aiRoutes.Post("/analyze-job", aiHandler.AnalyzeJob)

	// Export routes (protected)
	exportRoutes := api.Group("/export", r.authMiddleware.Protected())
	exportRoutes.Post("/pdf", exportHandler.ExportPDF)
	exportRoutes.Post("/txt", exportHandler.ExportTXT)
	exportRoutes.Post("/json", exportHandler.ExportJSON)

	return r.app
}

func (r *Router) Listen(addr string) error {
	return r.app.Listen(addr)
}
