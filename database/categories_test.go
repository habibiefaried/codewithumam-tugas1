package database

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/lib/pq"

	"codewithumam-tugas1/config"
)

// setupTestDB creates a test database connection and migrates the test table
func setupTestDB(t *testing.T) *sql.DB {
	// Check if secrets.yml exists - tests require it
	if _, err := os.Stat("secrets.yml"); os.IsNotExist(err) {
		t.Fatal("secrets.yml not found in current working directory - tests require database configuration from secrets.yml")
	}

	// Load config from secrets.yml (and env vars as fallback)
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	t.Logf("Loaded config - URL: %s, Port: %s, DB: %s, User: %s", cfg.URL, cfg.DBPort, cfg.Name, cfg.User)

	// Build connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.URL, cfg.DBPort, cfg.User, cfg.Password, cfg.Name)

	t.Logf("Connection string: host=%s port=%s user=%s password=%s dbname=%s", cfg.URL, cfg.DBPort, cfg.User, "***", cfg.Name)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}

	// Create test table
	if err := MigrateTest(db); err != nil {
		t.Fatalf("Failed to migrate test table: %v", err)
	}

	// Clean up any existing test data
	_, err = db.Exec("DELETE FROM category_test")
	if err != nil {
		t.Fatalf("Failed to clean up test data: %v", err)
	}

	return db
}

// teardownTestDB closes the database connection and drops the test table
func teardownTestDB(t *testing.T, db *sql.DB) {
	if err := DropTestTable(db); err != nil {
		t.Logf("Warning: failed to drop test table: %v", err)
	}

	if err := db.Close(); err != nil {
		t.Logf("Warning: failed to close database: %v", err)
	}
}

// TestCreateCategory tests the Create function
func TestCreateCategory(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	cat, err := Create(db, "category_test", "Electronics", "Electronic devices")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if cat.ID == 0 {
		t.Error("Expected non-zero ID")
	}

	if cat.Name != "Electronics" {
		t.Errorf("Expected name 'Electronics', got '%s'", cat.Name)
	}

	if cat.Description != "Electronic devices" {
		t.Errorf("Expected description 'Electronic devices', got '%s'", cat.Description)
	}
}

// TestGetByID tests the GetByID function
func TestGetByID(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	// Create a category first
	created, err := Create(db, "category_test", "Books", "All types of books")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Retrieve it
	cat, err := GetByID(db, "category_test", created.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if cat.ID != created.ID {
		t.Errorf("Expected ID %d, got %d", created.ID, cat.ID)
	}

	if cat.Name != "Books" {
		t.Errorf("Expected name 'Books', got '%s'", cat.Name)
	}

	if cat.Description != "All types of books" {
		t.Errorf("Expected description 'All types of books', got '%s'", cat.Description)
	}
}

// TestGetByIDNotFound tests GetByID when category doesn't exist
func TestGetByIDNotFound(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	_, err := GetByID(db, "category_test", 9999)
	if err == nil {
		t.Error("Expected error for non-existent category, got nil")
	}

	if err.Error() != "category not found" {
		t.Errorf("Expected 'category not found' error, got '%v'", err)
	}
}

// TestGetAll tests the GetAll function
func TestGetAll(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	// Create multiple categories
	_, err := Create(db, "category_test", "Electronics", "Electronic devices")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	_, err = Create(db, "category_test", "Books", "All types of books")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Get all
	categories, err := GetAll(db, "category_test")
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}

	if len(categories) != 2 {
		t.Errorf("Expected 2 categories, got %d", len(categories))
	}

	// Verify both categories exist
	names := map[string]bool{}
	for _, cat := range categories {
		names[cat.Name] = true
	}

	if !names["Electronics"] {
		t.Error("Expected 'Electronics' category not found")
	}

	if !names["Books"] {
		t.Error("Expected 'Books' category not found")
	}
}

// TestGetAllEmpty tests GetAll when no categories exist
func TestGetAllEmpty(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	categories, err := GetAll(db, "category_test")
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}

	if categories == nil {
		// This is expected - GetAll returns nil for empty result
		return
	}

	if len(categories) != 0 {
		t.Errorf("Expected 0 categories, got %d", len(categories))
	}
}

// TestUpdateCategory tests the Update function
func TestUpdateCategory(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	// Create a category
	created, err := Create(db, "category_test", "Electronics", "Electronic devices")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Update it
	updated, err := Update(db, "category_test", created.ID, "Updated Electronics", "Updated description")
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if updated.ID != created.ID {
		t.Errorf("Expected ID %d, got %d", created.ID, updated.ID)
	}

	if updated.Name != "Updated Electronics" {
		t.Errorf("Expected name 'Updated Electronics', got '%s'", updated.Name)
	}

	if updated.Description != "Updated description" {
		t.Errorf("Expected description 'Updated description', got '%s'", updated.Description)
	}

	// Verify by fetching
	fetched, err := GetByID(db, "category_test", created.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if fetched.Name != "Updated Electronics" {
		t.Errorf("Expected updated name 'Updated Electronics', got '%s'", fetched.Name)
	}
}

// TestUpdateCategoryNotFound tests Update when category doesn't exist
func TestUpdateCategoryNotFound(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	_, err := Update(db, "category_test", 9999, "Name", "Description")
	if err == nil {
		t.Error("Expected error for non-existent category, got nil")
	}

	if err.Error() != "category not found" {
		t.Errorf("Expected 'category not found' error, got '%v'", err)
	}
}

// TestDeleteCategory tests the Delete function
func TestDeleteCategory(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	// Create a category
	created, err := Create(db, "category_test", "Electronics", "Electronic devices")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Delete it
	err = Delete(db, "category_test", created.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify it's gone
	_, err = GetByID(db, "category_test", created.ID)
	if err == nil {
		t.Error("Expected category to be deleted, but still found")
	}
}

// TestDeleteCategoryNotFound tests Delete when category doesn't exist
func TestDeleteCategoryNotFound(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	err := Delete(db, "category_test", 9999)
	if err == nil {
		t.Error("Expected error for non-existent category, got nil")
	}

	if err.Error() != "category not found" {
		t.Errorf("Expected 'category not found' error, got '%v'", err)
	}
}

// TestMultipleOperations tests a sequence of operations
func TestMultipleOperations(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	// Create multiple categories
	cat1, err := Create(db, "category_test", "Electronics", "Electronic devices")
	if err != nil {
		t.Fatalf("Create cat1 failed: %v", err)
	}

	cat2, err := Create(db, "category_test", "Books", "All types of books")
	if err != nil {
		t.Fatalf("Create cat2 failed: %v", err)
	}

	// Get all
	categories, err := GetAll(db, "category_test")
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}

	if len(categories) != 2 {
		t.Errorf("Expected 2 categories, got %d", len(categories))
	}

	// Update one
	_, err = Update(db, "category_test", cat1.ID, "Updated Electronics", "Updated")
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Delete one
	err = Delete(db, "category_test", cat2.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Get all again
	categories, err = GetAll(db, "category_test")
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}

	if len(categories) != 1 {
		t.Errorf("Expected 1 category, got %d", len(categories))
	}

	if categories[0].Name != "Updated Electronics" {
		t.Errorf("Expected 'Updated Electronics', got '%s'", categories[0].Name)
	}
}
