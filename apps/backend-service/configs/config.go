package configs

import (
	"os"
)

type Config struct {
	Port               string
	DatabaseURL        string
	JWTSecret          string
	RefreshTokenSecret string
	FrontendURL        string
}

func LoadConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Default to local PostgreSQL docker container
		dbURL = "host=localhost user=admin password=admin dbname=ecommerce_db port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "ecommerce-super-secret-jwt-key-minimum-32-chars-long"
	}

	refreshTokenSecret := os.Getenv("REFRESH_TOKEN_SECRET")
	if refreshTokenSecret == "" {
		refreshTokenSecret = "ecommerce-super-secret-refresh-token-key-32-chars"
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}

	return &Config{
		Port:               port,
		DatabaseURL:        dbURL,
		JWTSecret:          jwtSecret,
		RefreshTokenSecret: refreshTokenSecret,
		FrontendURL:        frontendURL,
	}
}
