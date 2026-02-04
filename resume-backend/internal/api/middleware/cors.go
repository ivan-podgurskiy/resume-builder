package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/resume-builder/backend/internal/config"
)

func NewCORSMiddleware(cfg *config.Config) fiber.Handler {
	allowedOrigins := "http://localhost:5173,http://localhost:3000"
	if cfg.IsProduction() {
		allowedOrigins = "https://resume-builder.com,https://www.resume-builder.com"
	}

	return cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length,Content-Type,X-RateLimit-Limit,X-RateLimit-Remaining,X-RateLimit-Reset",
		MaxAge:           86400, // 24 hours
	})
}
