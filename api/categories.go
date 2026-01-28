package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"codewithumam-tugas1/database"
)

// Categories manages HTTP requests for categories
type Categories struct {
	db        *sql.DB
	tableName string
}

// NewCategories creates a new categories service
func NewCategories(db *sql.DB, tableName string) *Categories {
	return &Categories{
		db:        db,
		tableName: tableName,
	}
}

// GetAll handles GET /categories
func (c *Categories) GetAll(w http.ResponseWriter, r *http.Request) {
	categories, err := database.GetAll(c.db, c.tableName)
	if err != nil {
		http.Error(w, "Failed to retrieve categories", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if categories == nil {
		categories = []database.Category{}
	}
	json.NewEncoder(w).Encode(categories)
}

// GetByID handles GET /categories/{id}
func (c *Categories) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	cat, err := database.GetByID(c.db, c.tableName, id)
	if err != nil {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cat)
}

// Create handles POST /categories
func (c *Categories) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	cat, err := database.Create(c.db, c.tableName, req.Name, req.Description)
	if err != nil {
		http.Error(w, "Failed to create category", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cat)
}

// Update handles PUT /categories/{id}
func (c *Categories) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	cat, err := database.Update(c.db, c.tableName, id, req.Name, req.Description)
	if err != nil {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cat)
}

// Delete handles DELETE /categories/{id}
func (c *Categories) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = database.Delete(c.db, c.tableName, id)
	if err != nil {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
