package database

import (
	"database/sql"
	"fmt"
)

// Migrate creates the category table if it does not exist
func Migrate(db *sql.DB) error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS category (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		description TEXT
	);
	`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create category table: %w", err)
	}

	// Create indexes
	createIndexSQL := `
	CREATE INDEX IF NOT EXISTS idx_category_id ON category(id);
	CREATE INDEX IF NOT EXISTS idx_category_name ON category(name);
	`

	_, err = db.Exec(createIndexSQL)
	if err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	return nil
}

// MigrateTest creates the category_test table if it does not exist (for unit tests)
func MigrateTest(db *sql.DB) error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS category_test (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		description TEXT
	);
	`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create category_test table: %w", err)
	}

	// Create indexes
	createIndexSQL := `
	CREATE INDEX IF NOT EXISTS idx_category_test_id ON category_test(id);
	CREATE INDEX IF NOT EXISTS idx_category_test_name ON category_test(name);
	`

	_, err = db.Exec(createIndexSQL)
	if err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	return nil
}

// DropTestTable drops the category_test table (for cleanup in tests)
func DropTestTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS category_test;")
	if err != nil {
		return fmt.Errorf("failed to drop category_test table: %w", err)
	}
	return nil
}
