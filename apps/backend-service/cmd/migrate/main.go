package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"ecommerce/backend-service/configs"
	"ecommerce/backend-service/internal/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	cfg := configs.LoadConfig()

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("Fatal: Database connection failed: %v", err)
	}

	migrationDir := database.FindMigrationDir()

	cmd := "up"
	if len(os.Args) > 1 {
		cmd = strings.ToLower(os.Args[1])
	}

	switch cmd {
	case "up", "migrate":
		if err := database.MigrateUp(db, migrationDir); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
	case "down", "rollback":
		if err := database.MigrateRollback(db, migrationDir); err != nil {
			log.Fatalf("Rollback failed: %v", err)
		}
	case "fresh":
		withSeed := false
		for _, arg := range os.Args[2:] {
			if arg == "--seed" || arg == "-s" {
				withSeed = true
			}
		}
		if err := database.MigrateFresh(db, migrationDir, withSeed); err != nil {
			log.Fatalf("Fresh migration failed: %v", err)
		}
	case "seed":
		fmt.Println("Running database seeder...")
		database.SeedAllData(db, true)
		fmt.Println("Database seeded successfully.")
	case "status":
		if err := database.MigrateStatus(db, migrationDir); err != nil {
			log.Fatalf("Status check failed: %v", err)
		}
	default:
		fmt.Println("Usage: go run ./cmd/migrate [command]")
		fmt.Println("\nCommands (Laravel Artisan style):")
		fmt.Println("  up        Run all pending migrations (default)")
		fmt.Println("  rollback  Rollback the last applied migration")
		fmt.Println("  fresh     Drop all tables and re-run all migrations (--seed to also seed)")
		fmt.Println("  status    Show the status of each migration")
		fmt.Println("  seed      Run database seeder")
	}
}
