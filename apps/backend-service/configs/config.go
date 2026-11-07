package configs

import (
	"bufio"
	"os"
	"strings"
)

type Config struct {
	Port               string
	DatabaseURL        string
	JWTSecret          string
	RefreshTokenSecret string
	FrontendURL        string
}

func loadEnvFile(filenames ...string) {
	for _, filename := range filenames {
		file, err := os.Open(filename)
		if err != nil {
			continue
		}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				val := strings.TrimSpace(parts[1])
				// Don't overwrite if already set in environment
				if os.Getenv(key) == "" {
					os.Setenv(key, val)
				}
			}
		}
		_ = file.Close()
	}
}

func LoadConfig() *Config {
	// Try loading .env from current directory or app directory
	loadEnvFile(".env", "../../apps/backend-service/.env", "apps/backend-service/.env")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
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
