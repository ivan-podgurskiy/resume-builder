package middleware

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/resume-builder/backend/internal/config"
)

type RateLimiter struct {
	requests map[string]*rateLimitEntry
	mu       sync.RWMutex
	config   *config.Config
}

type rateLimitEntry struct {
	count     int
	resetTime time.Time
}

func NewRateLimiter(cfg *config.Config) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string]*rateLimitEntry),
		config:   cfg,
	}

	// Start cleanup goroutine
	go rl.cleanup()

	return rl
}

func (rl *RateLimiter) Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		key := c.IP()
		
		// Check if user is authenticated and use user ID as key
		if userID := GetUserID(c); userID.String() != "00000000-0000-0000-0000-000000000000" {
			key = userID.String()
		}

		rl.mu.Lock()
		entry, exists := rl.requests[key]
		now := time.Now()

		if !exists || now.After(entry.resetTime) {
			rl.requests[key] = &rateLimitEntry{
				count:     1,
				resetTime: now.Add(rl.config.RateLimitWindow),
			}
			rl.mu.Unlock()
			return c.Next()
		}

		if entry.count >= rl.config.RateLimitRequests {
			rl.mu.Unlock()
			
			c.Set("X-RateLimit-Limit", string(rune(rl.config.RateLimitRequests)))
			c.Set("X-RateLimit-Remaining", "0")
			c.Set("X-RateLimit-Reset", entry.resetTime.Format(time.RFC3339))
			c.Set("Retry-After", entry.resetTime.Sub(now).String())

			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":       "Rate limit exceeded",
				"retry_after": entry.resetTime.Sub(now).Seconds(),
			})
		}

		entry.count++
		remaining := rl.config.RateLimitRequests - entry.count
		rl.mu.Unlock()

		c.Set("X-RateLimit-Limit", string(rune(rl.config.RateLimitRequests)))
		c.Set("X-RateLimit-Remaining", string(rune(remaining)))
		c.Set("X-RateLimit-Reset", entry.resetTime.Format(time.RFC3339))

		return c.Next()
	}
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, entry := range rl.requests {
			if now.After(entry.resetTime) {
				delete(rl.requests, key)
			}
		}
		rl.mu.Unlock()
	}
}
