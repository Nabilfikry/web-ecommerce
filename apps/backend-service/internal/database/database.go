package database

import (
	"fmt"
	"log"
	"time"

	"ecommerce/backend-service/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// AutoMigrate fallback if database has not been migrated yet
	if !db.Migrator().HasTable(&model.User{}) {
		if err := db.AutoMigrate(
			&model.User{},
			&model.Product{},
			&model.Cart{},
			&model.CartItem{},
			&model.Order{},
			&model.OrderItem{},
			&model.RefreshToken{},
		); err != nil {
			return nil, fmt.Errorf("failed to auto-migrate database: %w", err)
		}
	}

	// Add unique index on LOWER(email) for case-insensitive email uniqueness
	_ = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_lower ON users (LOWER(email));").Error

	log.Println("Database connection established and schema migrated successfully")
	return db, nil
}

