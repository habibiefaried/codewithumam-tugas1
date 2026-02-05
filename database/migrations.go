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

	// Create product table
	createProductSQL := `
	CREATE TABLE IF NOT EXISTS product (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		price INTEGER NOT NULL DEFAULT 0,
		stock INTEGER NOT NULL DEFAULT 0,
		category_id INTEGER REFERENCES category(id)
	);
	`

	_, err = db.Exec(createProductSQL)
	if err != nil {
		return fmt.Errorf("failed to create product table: %w", err)
	}

	// Create product indexes
	createProductIndexSQL := `
	CREATE INDEX IF NOT EXISTS idx_product_id ON product(id);
	CREATE INDEX IF NOT EXISTS idx_product_name ON product(name);
	`

	_, err = db.Exec(createProductIndexSQL)
	if err != nil {
		return fmt.Errorf("failed to create product indexes: %w", err)
	}

	// Create transaction table
	createTransactionSQL := `
	CREATE TABLE IF NOT EXISTS "transaction" (
		id SERIAL PRIMARY KEY,
		total_amount INTEGER NOT NULL DEFAULT 0,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	`

	_, err = db.Exec(createTransactionSQL)
	if err != nil {
		return fmt.Errorf("failed to create transaction table: %w", err)
	}

	// Create transaction_detail table
	createTransactionDetailSQL := `
	CREATE TABLE IF NOT EXISTS transaction_detail (
		id SERIAL PRIMARY KEY,
		transaction_id INTEGER NOT NULL REFERENCES "transaction"(id) ON DELETE CASCADE,
		product_id INTEGER NOT NULL,
		product_name VARCHAR(255) NOT NULL,
		product_description TEXT NOT NULL DEFAULT '',
		unit_price INTEGER NOT NULL DEFAULT 0,
		quantity INTEGER NOT NULL,
		subtotal INTEGER NOT NULL
	);
	`

	_, err = db.Exec(createTransactionDetailSQL)
	if err != nil {
		return fmt.Errorf("failed to create transaction_detail table: %w", err)
	}

	alterTransactionDetailSQL := `
	ALTER TABLE transaction_detail DROP CONSTRAINT IF EXISTS transaction_detail_product_id_fkey;
	ALTER TABLE transaction_detail ADD COLUMN IF NOT EXISTS product_description TEXT NOT NULL DEFAULT '';
	ALTER TABLE transaction_detail ADD COLUMN IF NOT EXISTS unit_price INTEGER NOT NULL DEFAULT 0;
	`

	_, err = db.Exec(alterTransactionDetailSQL)
	if err != nil {
		return fmt.Errorf("failed to alter transaction_detail table: %w", err)
	}

	// Create transaction indexes
	createTransactionIndexSQL := `
	CREATE INDEX IF NOT EXISTS idx_transaction_id ON "transaction"(id);
	CREATE INDEX IF NOT EXISTS idx_transaction_created_at ON "transaction"(created_at);
	CREATE INDEX IF NOT EXISTS idx_transaction_created_at_id ON "transaction"(created_at, id);
	CREATE INDEX IF NOT EXISTS idx_transaction_detail_id ON transaction_detail(id);
	CREATE INDEX IF NOT EXISTS idx_transaction_detail_transaction_id ON transaction_detail(transaction_id);
	CREATE INDEX IF NOT EXISTS idx_transaction_detail_transaction_product_name ON transaction_detail(transaction_id, product_name);
	CREATE INDEX IF NOT EXISTS idx_transaction_detail_product_id ON transaction_detail(product_id);
	`

	_, err = db.Exec(createTransactionIndexSQL)
	if err != nil {
		return fmt.Errorf("failed to create transaction indexes: %w", err)
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

	// Create product_test table
	createProductSQL := `
	CREATE TABLE IF NOT EXISTS product_test (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		price INTEGER NOT NULL DEFAULT 0,
		stock INTEGER NOT NULL DEFAULT 0,
		category_id INTEGER
	);
	`

	_, err = db.Exec(createProductSQL)
	if err != nil {
		return fmt.Errorf("failed to create product_test table: %w", err)
	}

	// Create product_test indexes
	createProductIndexSQL := `
	CREATE INDEX IF NOT EXISTS idx_product_test_id ON product_test(id);
	CREATE INDEX IF NOT EXISTS idx_product_test_name ON product_test(name);
	`

	_, err = db.Exec(createProductIndexSQL)
	if err != nil {
		return fmt.Errorf("failed to create product_test indexes: %w", err)
	}

	// Create transaction_test table
	createTransactionTestSQL := `
	CREATE TABLE IF NOT EXISTS transaction_test (
		id SERIAL PRIMARY KEY,
		total_amount INTEGER NOT NULL DEFAULT 0,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	`

	_, err = db.Exec(createTransactionTestSQL)
	if err != nil {
		return fmt.Errorf("failed to create transaction_test table: %w", err)
	}

	// Create transaction_detail_test table
	createTransactionDetailTestSQL := `
	CREATE TABLE IF NOT EXISTS transaction_detail_test (
		id SERIAL PRIMARY KEY,
		transaction_id INTEGER NOT NULL REFERENCES transaction_test(id) ON DELETE CASCADE,
		product_id INTEGER NOT NULL,
		product_name VARCHAR(255) NOT NULL,
		product_description TEXT NOT NULL DEFAULT '',
		unit_price INTEGER NOT NULL DEFAULT 0,
		quantity INTEGER NOT NULL,
		subtotal INTEGER NOT NULL
	);
	`

	_, err = db.Exec(createTransactionDetailTestSQL)
	if err != nil {
		return fmt.Errorf("failed to create transaction_detail_test table: %w", err)
	}

	alterTransactionDetailTestSQL := `
	ALTER TABLE transaction_detail_test DROP CONSTRAINT IF EXISTS transaction_detail_test_product_id_fkey;
	ALTER TABLE transaction_detail_test ADD COLUMN IF NOT EXISTS product_description TEXT NOT NULL DEFAULT '';
	ALTER TABLE transaction_detail_test ADD COLUMN IF NOT EXISTS unit_price INTEGER NOT NULL DEFAULT 0;
	`

	_, err = db.Exec(alterTransactionDetailTestSQL)
	if err != nil {
		return fmt.Errorf("failed to alter transaction_detail_test table: %w", err)
	}

	// Create transaction_test indexes
	createTransactionTestIndexSQL := `
	CREATE INDEX IF NOT EXISTS idx_transaction_test_id ON transaction_test(id);
	CREATE INDEX IF NOT EXISTS idx_transaction_test_created_at ON transaction_test(created_at);
	CREATE INDEX IF NOT EXISTS idx_transaction_test_created_at_id ON transaction_test(created_at, id);
	CREATE INDEX IF NOT EXISTS idx_transaction_detail_test_id ON transaction_detail_test(id);
	CREATE INDEX IF NOT EXISTS idx_transaction_detail_test_transaction_id ON transaction_detail_test(transaction_id);
	CREATE INDEX IF NOT EXISTS idx_transaction_detail_test_transaction_product_name ON transaction_detail_test(transaction_id, product_name);
	CREATE INDEX IF NOT EXISTS idx_transaction_detail_test_product_id ON transaction_detail_test(product_id);
	`

	_, err = db.Exec(createTransactionTestIndexSQL)
	if err != nil {
		return fmt.Errorf("failed to create transaction_test indexes: %w", err)
	}

	return nil
}

// DropTestTable drops the category_test table (for cleanup in tests)
func DropTestTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS transaction_detail_test; DROP TABLE IF EXISTS transaction_test; DROP TABLE IF EXISTS product_test; DROP TABLE IF EXISTS category_test;")
	if err != nil {
		return fmt.Errorf("failed to drop category_test table: %w", err)
	}
	return nil
}
