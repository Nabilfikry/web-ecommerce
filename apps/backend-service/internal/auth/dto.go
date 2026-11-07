package auth

import (
	"time"

	"ecommerce/backend-service/internal/model"
	"github.com/google/uuid"
)

type RegisterRequest struct {
	Name     string         `json:"name"`
	Email    string         `json:"email"`
	Password string         `json:"password"`
	Role     model.UserRole `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID        uuid.UUID      `json:"id"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	Role      model.UserRole `json:"role"`
	CreatedAt time.Time      `json:"created_at"`
}

type LoginData struct {
	AccessToken string       `json:"access_token"`
	TokenType   string       `json:"token_type"`
	ExpiresIn   int          `json:"expires_in"` // seconds
	User        UserResponse `json:"user"`
}

type RefreshData struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}
