package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"codewithumam-tugas1/database"
)

// Checkout handles checkout requests
// It creates a transaction and updates product stock atomically.
type Checkout struct {
	db                     *sql.DB
	productTable           string
	transactionTable       string
	transactionDetailTable string
}

// NewCheckout creates a new checkout service
func NewCheckout(db *sql.DB, productTable, transactionTable, transactionDetailTable string) *Checkout {
	return &Checkout{
		db:                     db,
		productTable:           productTable,
		transactionTable:       transactionTable,
		transactionDetailTable: transactionDetailTable,
	}
}

// Create handles POST /checkout
func (c *Checkout) Create(w http.ResponseWriter, r *http.Request) {
	var req database.CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	transaction, err := database.Checkout(c.db, c.productTable, c.transactionTable, c.transactionDetailTable, req.Items)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrCheckoutEmptyItems), errors.Is(err, database.ErrInvalidCheckoutItem):
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		case errors.Is(err, database.ErrProductNotFound):
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		case errors.Is(err, database.ErrInsufficientStock):
			http.Error(w, "Insufficient stock", http.StatusBadRequest)
			return
		default:
			http.Error(w, "Failed to checkout", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transaction)
}
