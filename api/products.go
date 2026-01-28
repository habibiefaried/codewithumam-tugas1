package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"codewithumam-tugas1/database"
)

// Products manages HTTP requests for products
type Products struct {
	db        *sql.DB
	tableName string
}

// NewProducts creates a new products service
func NewProducts(db *sql.DB, tableName string) *Products {
	return &Products{db: db, tableName: tableName}
}

// GetAll handles GET /products
func (p *Products) GetAll(w http.ResponseWriter, r *http.Request) {
	products, err := database.GetAllProducts(p.db, p.tableName, "category")
	if err != nil {
		http.Error(w, "Failed to retrieve products", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if products == nil {
		products = []database.Product{}
	}
	json.NewEncoder(w).Encode(products)
}

// GetByID handles GET /products/{id}
func (p *Products) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	prod, err := database.GetProductByID(p.db, p.tableName, "category", id)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prod)
}

// Create handles POST /products
func (p *Products) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name       string `json:"name"`
		Price      int    `json:"price"`
		Stock      int    `json:"stock"`
		CategoryID int    `json:"category_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	prod, err := database.CreateProduct(p.db, p.tableName, "category", req.Name, req.Price, req.Stock, req.CategoryID)
	if err != nil {
		http.Error(w, "Failed to create product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(prod)
}

// Update handles PUT /products/{id}
func (p *Products) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Name       string `json:"name"`
		Price      int    `json:"price"`
		Stock      int    `json:"stock"`
		CategoryID int    `json:"category_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	prod, err := database.UpdateProduct(p.db, p.tableName, "category", id, req.Name, req.Price, req.Stock, req.CategoryID)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prod)
}

// Delete handles DELETE /products/{id}
func (p *Products) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = database.DeleteProduct(p.db, p.tableName, id)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
