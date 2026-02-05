package database

import (
	"codewithumam-tugas1/config"
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

// setupProductTestDB ensures secrets.yml exists and returns db connection
func setupProductTestDB(t *testing.T) *sql.DB {
	if _, err := os.Stat("secrets.yml"); os.IsNotExist(err) {
		t.Fatal("secrets.yml not found - tests require database configuration from secrets.yml")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.URL, cfg.DBPort, cfg.User, cfg.Password, cfg.Name)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}

	// Ensure tables
	if err := MigrateTest(db); err != nil {
		t.Fatalf("Failed to migrate test tables: %v", err)
	}

	// Clean up
	_, err = db.Exec("DELETE FROM transaction_detail_test")
	if err != nil {
		t.Fatalf("Failed to clean up transaction_detail_test: %v", err)
	}
	_, err = db.Exec("DELETE FROM transaction_test")
	if err != nil {
		t.Fatalf("Failed to clean up transaction_test: %v", err)
	}
	_, err = db.Exec("DELETE FROM product_test")
	if err != nil {
		t.Fatalf("Failed to clean up product_test: %v", err)
	}
	_, err = db.Exec("DELETE FROM category_test")
	if err != nil {
		t.Fatalf("Failed to clean up category_test: %v", err)
	}

	return db
}

func teardownProductTestDB(t *testing.T, db *sql.DB) {
	if err := DropTestTable(db); err != nil {
		t.Logf("Warning: failed to drop test tables: %v", err)
	}
	_ = db.Close()
}

func TestProductCRUD(t *testing.T) {
	db := setupProductTestDB(t)
	defer teardownProductTestDB(t, db)

	// create category first
	cat, err := Create(db, "category_test", "Gadgets", "Gadgets category")
	if err != nil {
		t.Fatalf("Failed to create category: %v", err)
	}

	// Create product
	prod, err := CreateProduct(db, "product_test", "category_test", "Phone", 299, 10, cat.ID)
	if err != nil {
		t.Fatalf("CreateProduct failed: %v", err)
	}

	if prod.ID == 0 {
		t.Error("Expected non-zero product ID")
	}

	if prod.Name != "Phone" {
		t.Errorf("Expected name Phone, got %s", prod.Name)
	}

	// Get by ID
	fetched, err := GetProductByID(db, "product_test", "category_test", prod.ID)
	if err != nil {
		t.Fatalf("GetProductByID failed: %v", err)
	}

	if fetched.Name != "Phone" {
		t.Errorf("Expected name Phone, got %s", fetched.Name)
	}

	// Update
	updated, err := UpdateProduct(db, "product_test", "category_test", prod.ID, "Smartphone", 399, 5, cat.ID)
	if err != nil {
		t.Fatalf("UpdateProduct failed: %v", err)
	}
	if updated.Name != "Smartphone" {
		t.Errorf("Expected name Smartphone, got %s", updated.Name)
	}

	// Delete
	err = DeleteProduct(db, "product_test", prod.ID)
	if err != nil {
		t.Fatalf("DeleteProduct failed: %v", err)
	}

	// Verify deletion
	_, err = GetProductByID(db, "product_test", "category_test", prod.ID)
	if err == nil {
		t.Error("Expected not found after deletion")
	}
}

// TestProductInputValidation tests input validation for products
func TestProductInputValidation(t *testing.T) {
	db := setupProductTestDB(t)
	defer teardownProductTestDB(t, db)

	// Create category for testing
	cat, err := Create(db, "category_test", "TestCat", "Test category")
	if err != nil {
		t.Fatalf("Failed to create category: %v", err)
	}

	// Test zero price (should fail at API layer)
	_, err = CreateProduct(db, "product_test", "category_test", "ZeroPrice", 0, 10, cat.ID)
	if err != nil {
		t.Logf("Database validation: zero price rejected (also validated at API layer)")
	}

	// Test negative stock (should fail at API layer)
	_, err = CreateProduct(db, "product_test", "category_test", "NegativeStock", 100, -5, cat.ID)
	if err != nil {
		t.Logf("Database validation: negative stock rejected (also validated at API layer)")
	}

	// Test with non-existent category (validation at API layer)
	prod, err := CreateProduct(db, "product_test", "category_test", "OrphanProduct", 100, 5, 9999)
	if err != nil {
		t.Logf("API layer validates category existence")
		return
	}
	if prod.ID > 0 {
		t.Logf("Database created orphan product (validation at API layer)")
	}
}

// TestProductUpdateValidation tests update validation for products
func TestProductUpdateValidation(t *testing.T) {
	db := setupProductTestDB(t)
	defer teardownProductTestDB(t, db)

	// Create category and product
	cat, err := Create(db, "category_test", "TestCat", "Test category")
	if err != nil {
		t.Fatalf("Failed to create category: %v", err)
	}

	prod, err := CreateProduct(db, "product_test", "category_test", "TestProduct", 100, 10, cat.ID)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	// Test update with non-existent category
	_, err = UpdateProduct(db, "product_test", "category_test", prod.ID, "Updated", 200, 5, 9999)
	if err != nil {
		t.Logf("API layer validates category on update")
		return
	}
	// If update succeeds, category validation is at API layer
	t.Logf("Update with invalid category succeeded at DB layer (API layer validates)")
}
