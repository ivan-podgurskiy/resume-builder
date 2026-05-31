package middleware

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/resume-builder/backend/internal/config"
)

func newTestMiddleware(secret string) *AuthMiddleware {
	return NewAuthMiddleware(&config.Config{
		JWTSecret:          secret,
		JWTExpirationHours: 24,
		RefreshTokenDays:   30,
	})
}

func TestGenerateAndValidateToken(t *testing.T) {
	m := newTestMiddleware("test-secret")
	userID := uuid.New()
	email := "user@example.com"

	token, err := m.GenerateToken(userID, email)
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}

	claims, err := m.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken returned error: %v", err)
	}
	if claims.UserID != userID {
		t.Errorf("UserID = %v, want %v", claims.UserID, userID)
	}
	if claims.Email != email {
		t.Errorf("Email = %q, want %q", claims.Email, email)
	}
	if claims.Issuer != "resume-builder" {
		t.Errorf("Issuer = %q, want %q", claims.Issuer, "resume-builder")
	}
}

func TestRefreshTokenHasDistinctIssuer(t *testing.T) {
	m := newTestMiddleware("test-secret")
	userID := uuid.New()

	token, err := m.GenerateRefreshToken(userID, "user@example.com")
	if err != nil {
		t.Fatalf("GenerateRefreshToken returned error: %v", err)
	}

	claims, err := m.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken returned error: %v", err)
	}
	if claims.Issuer != "resume-builder-refresh" {
		t.Errorf("Issuer = %q, want %q", claims.Issuer, "resume-builder-refresh")
	}
}

func TestValidateTokenRejectsWrongSecret(t *testing.T) {
	signer := newTestMiddleware("secret-a")
	verifier := newTestMiddleware("secret-b")

	token, err := signer.GenerateToken(uuid.New(), "user@example.com")
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}

	if _, err := verifier.ValidateToken(token); err == nil {
		t.Fatal("expected validation to fail for token signed with a different secret")
	}
}

func TestValidateTokenRejectsGarbage(t *testing.T) {
	m := newTestMiddleware("test-secret")
	if _, err := m.ValidateToken("not-a-jwt"); err == nil {
		t.Fatal("expected validation to fail for malformed token")
	}
}

func TestValidateTokenRejectsExpired(t *testing.T) {
	m := newTestMiddleware("test-secret")

	claims := JWTClaims{
		UserID: uuid.New(),
		Email:  "user@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			Issuer:    "resume-builder",
		},
	}
	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	if _, err := m.ValidateToken(signed); err == nil {
		t.Fatal("expected validation to fail for expired token")
	}
}
