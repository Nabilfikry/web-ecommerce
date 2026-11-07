package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"ecommerce/backend-service/internal/middleware"
	"ecommerce/backend-service/internal/shared"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux, authMiddleware func(http.Handler) http.Handler) {
	mux.HandleFunc("POST /api/v1/auth/register", h.handleRegister)
	mux.HandleFunc("POST /api/v1/auth/login", h.handleLogin)
	mux.HandleFunc("POST /api/v1/auth/refresh", h.handleRefresh)
	mux.HandleFunc("POST /api/v1/auth/logout", h.handleLogout)

	// Protected endpoint
	mux.Handle("GET /api/v1/auth/me", authMiddleware(http.HandlerFunc(h.handleMe)))
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.Error(w, http.StatusBadRequest, "Invalid request body JSON", "BAD_REQUEST", nil)
		return
	}

	user, err := h.service.Register(r.Context(), req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	shared.Success(w, http.StatusCreated, user, "User registered successfully")
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.Error(w, http.StatusBadRequest, "Invalid request body JSON", "BAD_REQUEST", nil)
		return
	}

	loginData, rawRefreshToken, err := h.service.Login(r.Context(), req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	// Set HttpOnly refresh token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    rawRefreshToken,
		Path:     "/api/v1/auth",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false, // set to true in production HTTPS
	})

	shared.Success(w, http.StatusOK, loginData, "Login successful")
}

func (h *Handler) handleRefresh(w http.ResponseWriter, r *http.Request) {
	var rawRefreshToken string

	// 1. Try reading from HttpOnly cookie
	if cookie, err := r.Cookie("refresh_token"); err == nil && cookie.Value != "" {
		rawRefreshToken = cookie.Value
	}

	// 2. Fallback to request body if sent
	if rawRefreshToken == "" && r.Body != nil {
		var body struct {
			RefreshToken string `json:"refresh_token"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		rawRefreshToken = body.RefreshToken
	}

	if rawRefreshToken == "" {
		shared.Error(w, http.StatusUnauthorized, "Refresh token is missing", "UNAUTHORIZED", nil)
		return
	}

	refreshData, err := h.service.Refresh(r.Context(), rawRefreshToken)
	if err != nil {
		h.handleError(w, err)
		return
	}

	shared.Success(w, http.StatusOK, refreshData, "Token refreshed successfully")
}

func (h *Handler) handleLogout(w http.ResponseWriter, r *http.Request) {
	var rawRefreshToken string
	if cookie, err := r.Cookie("refresh_token"); err == nil {
		rawRefreshToken = cookie.Value
	}

	_ = h.service.Logout(r.Context(), rawRefreshToken)

	// Invalidate cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/api/v1/auth",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	shared.Success(w, http.StatusOK, map[string]string{"message": "Logged out successfully"}, "Logged out successfully")
}

func (h *Handler) handleMe(w http.ResponseWriter, r *http.Request) {
	authUser, ok := middleware.GetAuthUser(r.Context())
	if !ok || authUser == nil {
		shared.Error(w, http.StatusUnauthorized, "User context not found", "UNAUTHORIZED", nil)
		return
	}

	user, err := h.service.GetMe(r.Context(), authUser.ID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	shared.Success(w, http.StatusOK, user, "")
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	var appErr *shared.AppError
	if errors.As(err, &appErr) {
		shared.Error(w, appErr.StatusCode, appErr.Message, appErr.Code, nil)
		return
	}
	shared.Error(w, http.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR", nil)
}
