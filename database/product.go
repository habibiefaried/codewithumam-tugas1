package database

// Product represents a product entity in the database. It includes
// CategoryName and CategoryDescription when joined with category table.
type Product struct {
	ID                  int    `json:"id" db:"id"`
	Name                string `json:"name" db:"name"`
	Price               int    `json:"price" db:"price"`
	Stock               int    `json:"stock" db:"stock"`
	CategoryID          int    `json:"category_id" db:"category_id"`
	CategoryName        string `json:"category_name" db:"category_name"`
	CategoryDescription string `json:"category_description" db:"category_description"`
}
