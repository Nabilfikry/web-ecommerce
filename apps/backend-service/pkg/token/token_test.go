package token

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestAccessTokenLifecycle(t *testing.T) {
	secret := "test-secret-key-that-is-sufficiently-long-32-bytes"
	userID := uuid.New()
	email := "buyer@example.com"
	role := "buyer"

	// 1. Valid token
	tokenStr, err := GenerateAccessToken(userID, email, role, secret, 15*time.Minute)
	if err != nil {
		t.Fatalf("failed to generate access token: %v", err)
	}

	claims, err := ValidateAccessToken(tokenStr, secret)
	if err != nil {
		t.Fatalf("failed to validate valid access token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("expected UserID %v, got %v", userID, claims.UserID)
	}
	if claims.Email != email {
		t.Errorf("expected Email %v, got %v", email, claims.Email)
	}
	if claims.Role != role {
		t.Errorf("expected Role %v, got %v", role, claims.Role)
	}

	// 2. Expired token
	expiredTokenStr, err := GenerateAccessToken(userID, email, role, secret, -1*time.Minute)
	if err != nil {
		t.Fatalf("failed to generate expired access token: %v", err)
	}

	_, err = ValidateAccessToken(expiredTokenStr, secret)
	if err == nil {
		t.Errorf("expected error when validating expired token, got nil")
	}

	// 3. Token with wrong secret
	_, err = ValidateAccessToken(tokenStr, "wrong-secret-key")
	if err == nil {
		t.Errorf("expected error when validating token with wrong secret, got nil")
	}
}

func TestRefreshTokenGeneration(t *testing.T) {
	raw, hash, err := GenerateRefreshToken()
	if err != nil {
		t.Fatalf("failed to generate refresh token: %v", err)
	}

	if raw == "" || hash == "" {
		t.Errorf("raw and hash tokens should not be empty")
	}

	computedHash := HashRefreshToken(raw)
	if computedHash != hash {
		t.Errorf("expected HashRefreshToken to match generated hash, got %v vs %v", computedHash, hash)
	}
}
