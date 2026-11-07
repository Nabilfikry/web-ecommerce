package database_test

import (
	"os"
	"path/filepath"
	"testing"

	"ecommerce/backend-service/configs"
	"ecommerce/backend-service/internal/database"
)

func TestMigrationWorkflow(t *testing.T) {
	cfg := configs.LoadConfig()
	db, err := database.InitDB(cfg.DatabaseURL)
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}

	// Verify EnsureMigrationTable
	if err := database.EnsureMigrationTable(db); err != nil {
		t.Fatalf("EnsureMigrationTable failed: %v", err)
	}

	// Verify FindMigrationDir
	dir := database.FindMigrationDir()
	if _, err := os.Stat(dir); err != nil {
		t.Fatalf("FindMigrationDir returned non-existent directory %q: %v", dir, err)
	}

	// Verify GetAppliedMigrations
	applied, err := database.GetAppliedMigrations(db)
	if err != nil {
		t.Fatalf("GetAppliedMigrations failed: %v", err)
	}

	if len(applied) == 0 {
		t.Log("Note: No migrations recorded yet in schema_migrations")
	}

	// Verify MigrateStatus executes without error
	if err := database.MigrateStatus(db, dir); err != nil {
		t.Fatalf("MigrateStatus failed: %v", err)
	}

	// Verify MigrateUp does not fail on already-applied migrations
	if err := database.MigrateUp(db, dir); err != nil {
		t.Fatalf("MigrateUp failed: %v", err)
	}
}

func TestFindMigrationDir_Fallback(t *testing.T) {
	origDir, _ := os.Getwd()
	tmpDir := t.TempDir()
	_ = os.Chdir(tmpDir)
	defer func() { _ = os.Chdir(origDir) }()

	dir := database.FindMigrationDir()
	// Should fallback to default string
	if filepath.Base(dir) != "migrations" {
		t.Errorf("expected fallback 'migrations', got %q", dir)
	}
}

