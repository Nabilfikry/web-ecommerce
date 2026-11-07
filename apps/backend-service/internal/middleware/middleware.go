package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"ecommerce/backend-service/internal/shared"
	"ecommerce/backend-service/pkg/token"
	"github.com/google/uuid"
)

type contextKey string

const (
	UserContextKey contextKey = "authUser"
)

type AuthUser struct {
	ID    uuid.UUID
	Email string
	Role  string
}

func GetAuthUser(ctx context.Context) (*AuthUser, bool) {
	u, ok := ctx.Value(UserContextKey).(*AuthUser)
	return u, ok
}

// CORS Middleware handles cross-origin requests and preflight
func CORS(allowedOrigin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin == "" {
				origin = allowedOrigin
			}

			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "86400")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// AuthMiddleware validates the JWT access token from Authorization header
func Auth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				shared.Error(w, http.StatusUnauthorized, "Missing Authorization header", "UNAUTHORIZED", nil)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				shared.Error(w, http.StatusUnauthorized, "Invalid Authorization header format. Expected 'Bearer <token>'", "UNAUTHORIZED", nil)
				return
			}

			tokenStr := parts[1]
			claims, err := token.ValidateAccessToken(tokenStr, jwtSecret)
			if err != nil {
				shared.Error(w, http.StatusUnauthorized, "Invalid or expired access token", "UNAUTHORIZED", nil)
				return
			}

			authUser := &AuthUser{
				ID:    claims.UserID,
				Email: claims.Email,
				Role:  claims.Role,
			}

			ctx := context.WithValue(r.Context(), UserContextKey, authUser)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole checks if the authenticated user has one of the allowed roles
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := GetAuthUser(r.Context())
			if !ok || user == nil {
				shared.Error(w, http.StatusUnauthorized, "User context not found", "UNAUTHORIZED", nil)
				return
			}

			allowed := false
			for _, rRole := range roles {
				if user.Role == rRole {
					allowed = true
					break
				}
			}

			if !allowed {
				shared.Error(w, http.StatusForbidden, "Access forbidden for current user role", "FORBIDDEN", nil)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Recovery recovers from panics and returns 500
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("[PANIC RECOVERED] %v\n", rec)
				shared.Error(w, http.StatusInternalServerError, "Internal server error occurred", "INTERNAL_ERROR", nil)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

// Logger logs each incoming HTTP request
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rec, r)
		log.Printf("%s %s %d %v\n", r.Method, r.URL.Path, rec.statusCode, time.Since(start))
	})
}
