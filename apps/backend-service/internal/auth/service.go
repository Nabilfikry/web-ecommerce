package auth

import (
	"context"
	"errors"
	"net/mail"
	"strings"
	"time"

	"ecommerce/backend-service/internal/model"
	"ecommerce/backend-service/internal/shared"
	"ecommerce/backend-service/pkg/hash"
	"ecommerce/backend-service/pkg/token"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	Register(ctx context.Context, req RegisterRequest) (*UserResponse, error)
	Login(ctx context.Context, req LoginRequest) (*LoginData, string, error)
	Refresh(ctx context.Context, rawRefreshToken string) (*RefreshData, error)
	Logout(ctx context.Context, rawRefreshToken string) error
	GetMe(ctx context.Context, userID uuid.UUID) (*UserResponse, error)
}

type authService struct {
	repo               Repository
	jwtSecret          string
	refreshTokenSecret string
}

func NewService(repo Repository, jwtSecret, refreshTokenSecret string) Service {
	return &authService{
		repo:               repo,
		jwtSecret:          jwtSecret,
		refreshTokenSecret: refreshTokenSecret,
	}
}

func (s *authService) Register(ctx context.Context, req RegisterRequest) (*UserResponse, error) {
	name := strings.TrimSpace(req.Name)
	if len(name) < 2 || len(name) > 100 {
		return nil, shared.NewAppError("VALIDATION_ERROR", "Name must be between 2 and 100 characters", 400, nil)
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, shared.NewAppError("VALIDATION_ERROR", "Invalid email format", 400, nil)
	}

	if len(req.Password) < 8 {
		return nil, shared.NewAppError("VALIDATION_ERROR", "Password must be at least 8 characters", 400, nil)
	}

	if req.Role != model.RoleBuyer && req.Role != model.RoleSeller {
		return nil, shared.NewAppError("VALIDATION_ERROR", "Role must be 'buyer' or 'seller'", 400, nil)
	}

	// Check existing user
	existing, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, shared.NewAppError("EMAIL_ALREADY_REGISTERED", "Email is already registered", 409, nil)
	}

	hashedPassword, err := hash.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:           uuid.New(),
		Name:         name,
		Email:        email,
		PasswordHash: hashedPassword,
		Role:         req.Role,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	// If role is buyer, auto-create a shopping cart
	if req.Role == model.RoleBuyer {
		cart := &model.Cart{
			ID:      uuid.New(),
			BuyerID: user.ID,
		}
		_ = s.repo.CreateCart(ctx, cart)
	}

	return &UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *authService) Login(ctx context.Context, req LoginRequest) (*LoginData, string, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", shared.NewAppError("INVALID_CREDENTIALS", "Invalid email or password", 401, nil)
		}
		return nil, "", err
	}

	if !hash.CheckPassword(req.Password, user.PasswordHash) {
		return nil, "", shared.NewAppError("INVALID_CREDENTIALS", "Invalid email or password", 401, nil)
	}

	// Generate access token (15 mins = 900 seconds)
	accessToken, err := token.GenerateAccessToken(user.ID, user.Email, string(user.Role), s.jwtSecret, 15*time.Minute)
	if err != nil {
		return nil, "", err
	}

	// Generate refresh token (7 days)
	rawRefreshToken, tokenHash, err := token.GenerateRefreshToken()
	if err != nil {
		return nil, "", err
	}

	refreshTokenModel := &model.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	if err := s.repo.StoreRefreshToken(ctx, refreshTokenModel); err != nil {
		return nil, "", err
	}

	return &LoginData{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   900,
		User: UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		},
	}, rawRefreshToken, nil
}

func (s *authService) Refresh(ctx context.Context, rawRefreshToken string) (*RefreshData, error) {
	if rawRefreshToken == "" {
		return nil, shared.NewAppError("INVALID_TOKEN", "Refresh token is missing", 401, nil)
	}

	tokenHash := token.HashRefreshToken(rawRefreshToken)
	rfToken, err := s.repo.GetRefreshToken(ctx, tokenHash)
	if err != nil {
		return nil, shared.NewAppError("INVALID_TOKEN", "Invalid or expired refresh token", 401, nil)
	}

	if rfToken.RevokedAt != nil || rfToken.ExpiresAt.Before(time.Now()) {
		return nil, shared.NewAppError("INVALID_TOKEN", "Refresh token is revoked or expired", 401, nil)
	}

	user, err := s.repo.GetUserByID(ctx, rfToken.UserID)
	if err != nil {
		return nil, shared.NewAppError("USER_NOT_FOUND", "User associated with token not found", 401, nil)
	}

	accessToken, err := token.GenerateAccessToken(user.ID, user.Email, string(user.Role), s.jwtSecret, 15*time.Minute)
	if err != nil {
		return nil, err
	}

	return &RefreshData{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   900,
	}, nil
}

func (s *authService) Logout(ctx context.Context, rawRefreshToken string) error {
	if rawRefreshToken == "" {
		return nil
	}
	tokenHash := token.HashRefreshToken(rawRefreshToken)
	return s.repo.RevokeRefreshToken(ctx, tokenHash)
}

func (s *authService) GetMe(ctx context.Context, userID uuid.UUID) (*UserResponse, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, shared.NewAppError("USER_NOT_FOUND", "User not found", 404, nil)
	}

	return &UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}, nil
}

