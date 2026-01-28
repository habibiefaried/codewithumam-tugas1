package database

import (
	"database/sql"
	"fmt"
)

// GetAll retrieves all categories from the database
func GetAll(db *sql.DB, tableName string) ([]Category, error) {
	query := fmt.Sprintf("SELECT id, name, description FROM %s", tableName)
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query categories: %w", err)
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var cat Category
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.Description); err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, cat)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating categories: %w", err)
	}

	return categories, nil
}

// GetByID retrieves a category by ID from the database
func GetByID(db *sql.DB, tableName string, id int) (Category, error) {
	var cat Category
	query := fmt.Sprintf("SELECT id, name, description FROM %s WHERE id = $1", tableName)
	err := db.QueryRow(query, id).Scan(&cat.ID, &cat.Name, &cat.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return Category{}, fmt.Errorf("category not found")
		}
		return Category{}, fmt.Errorf("failed to query category: %w", err)
	}
	return cat, nil
}

// Create inserts a new category into the database and returns the created category
func Create(db *sql.DB, tableName string, name, description string) (Category, error) {
	var cat Category
	query := fmt.Sprintf("INSERT INTO %s (name, description) VALUES ($1, $2) RETURNING id, name, description", tableName)
	err := db.QueryRow(query, name, description).Scan(&cat.ID, &cat.Name, &cat.Description)
	if err != nil {
		return Category{}, fmt.Errorf("failed to create category: %w", err)
	}
	return cat, nil
}

// Update modifies an existing category in the database
func Update(db *sql.DB, tableName string, id int, name, description string) (Category, error) {
	var cat Category
	query := fmt.Sprintf("UPDATE %s SET name = $1, description = $2 WHERE id = $3 RETURNING id, name, description", tableName)
	err := db.QueryRow(query, name, description, id).Scan(&cat.ID, &cat.Name, &cat.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return Category{}, fmt.Errorf("category not found")
		}
		return Category{}, fmt.Errorf("failed to update category: %w", err)
	}
	return cat, nil
}

// Delete removes a category from the database
func Delete(db *sql.DB, tableName string, id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", tableName)
	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("category not found")
	}

	return nil
}
