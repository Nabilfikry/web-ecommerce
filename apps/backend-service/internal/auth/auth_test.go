package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"ecommerce/backend-service/internal/database"
	"ecommerce/backend-service/internal/middleware"
	"ecommerce/backend-service/internal/model"
	"ecommerce/backend-service/pkg/token"
	"github.com/google/uuid"
)

func setupTestDB(t *testing.T) Repository {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=admin password=admin dbname=ecommerce_db port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	}
	db, err := database.InitDB(dsn)
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}
	return NewRepository(db)
}

func TestAuthService_RegisterAndLogin(t *testing.T) {
	repo := setupTestDB(t)
	jwtSecret := "test-jwt-secret-minimum-32-chars-long"
	rfSecret := "test-rf-secret-minimum-32-chars-long"
	svc := NewService(repo, jwtSecret, rfSecret)
	ctx := context.Background()

	randomSuffix := uuid.New().String()[:8]
	testEmail := "user_" + randomSuffix + "@example.com"

	// 1. Invalid registration inputs
	_, err := svc.Register(ctx, RegisterRequest{
		Name:     "A",
		Email:    testEmail,
		Password: "Password123",
		Role:     model.RoleBuyer,
	})
	if err == nil {
		t.Errorf("expected error for short name, got nil")
	}

	_, err = svc.Register(ctx, RegisterRequest{
		Name:     "Valid Name",
		Email:    "invalid-email",
		Password: "Password123",
		Role:     model.RoleBuyer,
	})
	if err == nil {
		t.Errorf("expected error for invalid email format, got nil")
	}

	_, err = svc.Register(ctx, RegisterRequest{
		Name:     "Valid Name",
		Email:    testEmail,
		Password: "123",
		Role:     model.RoleBuyer,
	})
	if err == nil {
		t.Errorf("expected error for short password, got nil")
	}

	// 2. Successful Registration
	regRes, err := svc.Register(ctx, RegisterRequest{
		Name:     "Test Buyer",
		Email:    testEmail,
		Password: "SecurePassword123",
		Role:     model.RoleBuyer,
	})
	if err != nil {
		t.Fatalf("failed to register valid user: %v", err)
	}

	if regRes.Email != testEmail {
		t.Errorf("expected email %s, got %s", testEmail, regRes.Email)
	}
	if regRes.Role != model.RoleBuyer {
		t.Errorf("expected role buyer, got %s", regRes.Role)
	}

	// 3. Duplicate Email Registration
	_, err = svc.Register(ctx, RegisterRequest{
		Name:     "Another Name",
		Email:    testEmail,
		Password: "SecurePassword123",
		Role:     model.RoleBuyer,
	})
	if err == nil {
		t.Errorf("expected conflict error for duplicate email, got nil")
	}

	// 4. Login with Wrong Password
	_, _, err = svc.Login(ctx, LoginRequest{
		Email:    testEmail,
		Password: "WrongPassword!",
	})
	if err == nil {
		t.Errorf("expected error for wrong password, got nil")
	}

	// 5. Successful Login
	loginRes, rfToken, err := svc.Login(ctx, LoginRequest{
		Email:    testEmail,
		Password: "SecurePassword123",
	})
	if err != nil {
		t.Fatalf("failed to login: %v", err)
	}

	if loginRes.AccessToken == "" || rfToken == "" {
		t.Errorf("access token and refresh token must not be empty")
	}

	// 6. Refresh Token
	refreshed, err := svc.Refresh(ctx, rfToken)
	if err != nil {
		t.Fatalf("failed to refresh token: %v", err)
	}
	if refreshed.AccessToken == "" {
		t.Errorf("expected new access token, got empty")
	}

	// 7. GetMe
	me, err := svc.GetMe(ctx, regRes.ID)
	if err != nil {
		t.Fatalf("failed to get me: %v", err)
	}
	if me.Name != "Test Buyer" {
		t.Errorf("expected user name 'Test Buyer', got %s", me.Name)
	}

	// 8. Logout & Revoke
	if err := svc.Logout(ctx, rfToken); err != nil {
		t.Fatalf("failed to logout: %v", err)
	}

	// 9. Re-attempt refresh with revoked token
	_, err = svc.Refresh(ctx, rfToken)
	if err == nil {
		t.Errorf("expected error refreshing with revoked token, got nil")
	}
}

func TestAuthMiddleware_And_RequireRole(t *testing.T) {
	jwtSecret := "test-jwt-secret-minimum-32-chars-long"

	// Create test handler protected by auth and role middleware
	protectedHandler := middleware.Auth(jwtSecret)(
		middleware.RequireRole("seller")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("seller authorized"))
		})),
	)

	// 1. Request without header -> 401
	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()
	protectedHandler.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for missing header, got %d", rec.Code)
	}

	// 2. Request with buyer token to seller route -> 403
	buyerToken, err := token.GenerateAccessToken(uuid.New(), "buyer@test.com", "buyer", jwtSecret, 15*time.Minute)
	if err != nil {
		t.Fatalf("failed to generate buyer token: %v", err)
	}

	req2 := httptest.NewRequest("GET", "/test", nil)
	req2.Header.Set("Authorization", "Bearer "+buyerToken)
	rec2 := httptest.NewRecorder()
	protectedHandler.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusForbidden {
		t.Errorf("expected 403 for buyer accessing seller route, got %d", rec2.Code)
	}

	// 3. Request with seller token to seller route -> 200
	sellerToken, err := token.GenerateAccessToken(uuid.New(), "seller@test.com", "seller", jwtSecret, 15*time.Minute)
	if err != nil {
		t.Fatalf("failed to generate seller token: %v", err)
	}

	req3 := httptest.NewRequest("GET", "/test", nil)
	req3.Header.Set("Authorization", "Bearer "+sellerToken)
	rec3 := httptest.NewRecorder()
	protectedHandler.ServeHTTP(rec3, req3)
	if rec3.Code != http.StatusOK {
		t.Errorf("expected 200 for seller accessing seller route, got %d", rec3.Code)
	}
}
