package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

// GitCommit will be set at build time
var GitCommit string

// Category represents a category entity
type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CategoryStore handles in-memory storage for categories
type CategoryStore struct {
	mu         sync.RWMutex
	categories map[int]Category
	nextID     int
}

// NewCategoryStore creates a new category store
func NewCategoryStore() *CategoryStore {
	return &CategoryStore{
		categories: make(map[int]Category),
		nextID:     1,
	}
}

// GetAll returns all categories
func (s *CategoryStore) GetAll() []Category {
	s.mu.RLock()
	defer s.mu.RUnlock()

	categories := make([]Category, 0, len(s.categories))
	for _, cat := range s.categories {
		categories = append(categories, cat)
	}
	return categories
}

// GetByID returns a category by ID
func (s *CategoryStore) GetByID(id int) (Category, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cat, exists := s.categories[id]
	return cat, exists
}

// Create adds a new category
func (s *CategoryStore) Create(name, description string) Category {
	s.mu.Lock()
	defer s.mu.Unlock()

	cat := Category{
		ID:          s.nextID,
		Name:        name,
		Description: description,
	}
	s.categories[s.nextID] = cat
	s.nextID++
	return cat
}

// Update modifies an existing category
func (s *CategoryStore) Update(id int, name, description string) (Category, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.categories[id]; !exists {
		return Category{}, false
	}

	cat := Category{
		ID:          id,
		Name:        name,
		Description: description,
	}
	s.categories[id] = cat
	return cat, true
}

// Delete removes a category
func (s *CategoryStore) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.categories[id]; !exists {
		return false
	}

	delete(s.categories, id)
	return true
}

// CategoryHandler handles HTTP requests for categories
type CategoryHandler struct {
	store *CategoryStore
}

// HandleGetAll handles GET /categories
func (h *CategoryHandler) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	categories := h.store.GetAll()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// HandleGetByID handles GET /categories/{id}
func (h *CategoryHandler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	cat, exists := h.store.GetByID(id)
	if !exists {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cat)
}

// HandleCreate handles POST /categories
func (h *CategoryHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
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

	cat := h.store.Create(req.Name, req.Description)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cat)
}

// HandleUpdate handles PUT /categories/{id}
func (h *CategoryHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
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

	cat, exists := h.store.Update(id, req.Name, req.Description)
	if !exists {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cat)
}

// HandleDelete handles DELETE /categories/{id}
func (h *CategoryHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !h.store.Delete(id) {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize category store and handler
	store := NewCategoryStore()
	handler := &CategoryHandler{store: store}

	// Category routes
	http.HandleFunc("GET /categories", handler.HandleGetAll)
	http.HandleFunc("GET /categories/{id}", handler.HandleGetByID)
	http.HandleFunc("POST /categories", handler.HandleCreate)
	http.HandleFunc("PUT /categories/{id}", handler.HandleUpdate)
	http.HandleFunc("DELETE /categories/{id}", handler.HandleDelete)

	// Original endpoints
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!")
	})

	http.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})

	http.HandleFunc("GET /version", func(w http.ResponseWriter, r *http.Request) {
		if GitCommit == "" {
			GitCommit = "unknown"
		}
		fmt.Fprintf(w, "Commit: %s", GitCommit)
	})

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(addr, nil))
}
