module github.com/resume-builder/backend

go 1.24

require (
	github.com/gofiber/fiber/v2 v2.52.0
	github.com/gofiber/contrib/jwt v1.0.8
	github.com/gofiber/contrib/websocket v1.3.0
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/google/uuid v1.6.0
	github.com/hibiken/asynq v0.24.1
	github.com/joho/godotenv v1.5.1
	github.com/rs/zerolog v1.31.0
	github.com/go-playground/validator/v10 v10.17.0
	github.com/anthropics/anthropic-sdk-go v0.2.0-alpha.4
	gorm.io/gorm v1.25.6
	gorm.io/driver/postgres v1.5.4
	golang.org/x/crypto v0.18.0
	github.com/aws/aws-sdk-go-v2 v1.24.1
	github.com/aws/aws-sdk-go-v2/config v1.26.6
	github.com/aws/aws-sdk-go-v2/service/s3 v1.48.1
	github.com/aws/aws-sdk-go-v2/credentials v1.16.16
)
