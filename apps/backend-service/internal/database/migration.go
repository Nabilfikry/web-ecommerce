package database

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"
)

type MigrationRecord struct {
	Version   string    `gorm:"primaryKey;type:varchar(150)"`
	AppliedAt time.Time `gorm:"not null"`
}

func EnsureMigrationTable(db *gorm.DB) error {
	return db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(150) PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`).Error
}

func GetAppliedMigrations(db *gorm.DB) (map[string]time.Time, error) {
	applied := make(map[string]time.Time)
	var records []MigrationRecord
	if err := db.Table("schema_migrations").Find(&records).Error; err != nil {
		return nil, err
	}
	for _, r := range records {
		applied[r.Version] = r.AppliedAt
	}
	return applied, nil
}

func FindMigrationDir() string {
	candidates := []string{
		"migrations",
		"../migrations",
		"../../migrations",
		"apps/backend-service/migrations",
		"../../apps/backend-service/migrations",
	}
	for _, c := range candidates {
		if fi, err := os.Stat(c); err == nil && fi.IsDir() {
			return c
		}
	}
	return "migrations"
}

func MigrateUp(db *gorm.DB, dir string) error {
	if err := EnsureMigrationTable(db); err != nil {
		return err
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read migrations dir: %w", err)
	}

	applied, err := GetAppliedMigrations(db)
	if err != nil {
		return err
	}

	var upFiles []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".up.sql") {
			upFiles = append(upFiles, f.Name())
		}
	}
	sort.Strings(upFiles)

	count := 0
	for _, filename := range upFiles {
		version := strings.TrimSuffix(filename, ".up.sql")
		if _, exists := applied[version]; exists {
			continue
		}

		filePath := filepath.Join(dir, filename)
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", filename, err)
		}

		start := time.Now()
		fmt.Printf("Migrating: %s\n", filename)

		err = db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Exec(string(content)).Error; err != nil {
				return err
			}
			return tx.Table("schema_migrations").Create(&MigrationRecord{
				Version:   version,
				AppliedAt: time.Now(),
			}).Error
		})

		if err != nil {
			return fmt.Errorf("failed executing %s: %w", filename, err)
		}

		fmt.Printf("Migrated:  %s (%v)\n", filename, time.Since(start).Round(time.Millisecond))
		count++
	}

	if count == 0 {
		fmt.Println("Nothing to migrate. Database is already up to date.")
	} else {
		fmt.Printf("Successfully applied %d migration(s).\n", count)
	}
	return nil
}

func MigrateRollback(db *gorm.DB, dir string) error {
	if err := EnsureMigrationTable(db); err != nil {
		return err
	}

	var latest MigrationRecord
	err := db.Table("schema_migrations").Order("applied_at DESC, version DESC").First(&latest).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Println("No migrations to rollback.")
			return nil
		}
		return err
	}

	downFilename := latest.Version + ".down.sql"
	filePath := filepath.Join(dir, downFilename)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read rollback file %s: %w", downFilename, err)
	}

	start := time.Now()
	fmt.Printf("Rolling back: %s\n", downFilename)

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(string(content)).Error; err != nil {
			return err
		}
		return tx.Table("schema_migrations").Where("version = ?", latest.Version).Delete(&MigrationRecord{}).Error
	})

	if err != nil {
		return fmt.Errorf("failed rolling back %s: %w", downFilename, err)
	}

	fmt.Printf("Rolled back:  %s (%v)\n", downFilename, time.Since(start).Round(time.Millisecond))
	return nil
}

func MigrateFresh(db *gorm.DB, dir string, withSeed bool) error {
	fmt.Println("Dropping all tables from database...")

	dropSQL := `
		DROP SCHEMA public CASCADE;
		CREATE SCHEMA public;
		GRANT ALL ON SCHEMA public TO public;
	`
	if err := db.Exec(dropSQL).Error; err != nil {
		return fmt.Errorf("failed to drop schema: %w", err)
	}

	fmt.Println("Schema reset. Re-running all migrations...")
	if err := EnsureMigrationTable(db); err != nil {
		return err
	}

	if err := MigrateUp(db, dir); err != nil {
		return err
	}

	if withSeed {
		fmt.Println("Running database seeder...")
		SeedInitialData(db)
		fmt.Println("Database seeded successfully.")
	}

	return nil
}

func MigrateStatus(db *gorm.DB, dir string) error {
	if err := EnsureMigrationTable(db); err != nil {
		return err
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read migrations dir: %w", err)
	}

	applied, err := GetAppliedMigrations(db)
	if err != nil {
		return err
	}

	var upFiles []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".up.sql") {
			upFiles = append(upFiles, f.Name())
		}
	}
	sort.Strings(upFiles)

	fmt.Printf("\n%-40s | %-10s | %-25s\n", "Migration", "Status", "Applied At (UTC)")
	fmt.Println(strings.Repeat("-", 80))

	for _, filename := range upFiles {
		version := strings.TrimSuffix(filename, ".up.sql")
		if appliedAt, ok := applied[version]; ok {
			fmt.Printf("%-40s | \033[32m%-10s\033[0m | %-25s\n", version, "Ran", appliedAt.Format(time.RFC3339))
		} else {
			fmt.Printf("%-40s | \033[33m%-10s\033[0m | %-25s\n", version, "Pending", "-")
		}
	}
	fmt.Println()
	return nil
}
