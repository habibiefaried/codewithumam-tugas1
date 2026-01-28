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

// Product queries

// GetAllProducts retrieves all products with category info
func GetAllProducts(db *sql.DB, tableName, categoryTableName string) ([]Product, error) {
	query := fmt.Sprintf(`SELECT p.id, p.name, p.price, p.stock, p.category_id, COALESCE(c.name, ''), COALESCE(c.description, '') FROM %s p LEFT JOIN %s c ON p.category_id = c.id`, tableName, categoryTableName)
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query products: %w", err)
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID, &p.CategoryName, &p.CategoryDescription); err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating products: %w", err)
	}

	return products, nil
}

// GetProductByID retrieves a product by ID with category info
func GetProductByID(db *sql.DB, tableName, categoryTableName string, id int) (Product, error) {
	var p Product
	query := fmt.Sprintf(`SELECT p.id, p.name, p.price, p.stock, p.category_id, COALESCE(c.name, ''), COALESCE(c.description, '') FROM %s p LEFT JOIN %s c ON p.category_id = c.id WHERE p.id = $1`, tableName, categoryTableName)
	err := db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID, &p.CategoryName, &p.CategoryDescription)
	if err != nil {
		if err == sql.ErrNoRows {
			return Product{}, fmt.Errorf("product not found")
		}
		return Product{}, fmt.Errorf("failed to query product: %w", err)
	}
	return p, nil
}

// CreateProduct inserts a new product into the database and returns it
func CreateProduct(db *sql.DB, tableName, categoryTableName string, name string, price, stock, categoryID int) (Product, error) {
	var p Product
	query := fmt.Sprintf("INSERT INTO %s (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id, name, price, stock, category_id", tableName)
	err := db.QueryRow(query, name, price, stock, categoryID).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID)
	if err != nil {
		return Product{}, fmt.Errorf("failed to create product: %w", err)
	}
	// load category info
	if p.CategoryID != 0 {
		cat, _ := GetByID(db, categoryTableName, p.CategoryID)
		p.CategoryName = cat.Name
		p.CategoryDescription = cat.Description
	}
	return p, nil
}

// UpdateProduct updates an existing product and returns it
func UpdateProduct(db *sql.DB, tableName, categoryTableName string, id int, name string, price, stock, categoryID int) (Product, error) {
	var p Product
	query := fmt.Sprintf("UPDATE %s SET name = $1, price = $2, stock = $3, category_id = $4 WHERE id = $5 RETURNING id, name, price, stock, category_id", tableName)
	err := db.QueryRow(query, name, price, stock, categoryID, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return Product{}, fmt.Errorf("product not found")
		}
		return Product{}, fmt.Errorf("failed to update product: %w", err)
	}
	if p.CategoryID != 0 {
		cat, _ := GetByID(db, categoryTableName, p.CategoryID)
		p.CategoryName = cat.Name
		p.CategoryDescription = cat.Description
	}
	return p, nil
}

// DeleteProduct removes a product
func DeleteProduct(db *sql.DB, tableName string, id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", tableName)
	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("product not found")
	}
	return nil
}
