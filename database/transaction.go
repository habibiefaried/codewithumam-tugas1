package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var (
	ErrCheckoutEmptyItems  = errors.New("checkout items required")
	ErrInvalidCheckoutItem = errors.New("invalid checkout item")
	ErrProductNotFound     = errors.New("product not found")
	ErrInsufficientStock   = errors.New("insufficient stock")
)

// Transaction represents a checkout transaction with details
// It includes a timestamp for reporting.
type Transaction struct {
	ID          int                 `json:"id" db:"id"`
	TotalAmount int                 `json:"total_amount" db:"total_amount"`
	CreatedAt   time.Time           `json:"created_at" db:"created_at"`
	Details     []TransactionDetail `json:"details" db:"-"`
}

// TransactionDetail represents each item in a transaction
type TransactionDetail struct {
	ID            int    `json:"id" db:"id"`
	TransactionID int    `json:"transaction_id" db:"transaction_id"`
	ProductID     int    `json:"product_id" db:"product_id"`
	ProductName   string `json:"product_name" db:"product_name"`
	ProductDesc   string `json:"product_description" db:"product_description"`
	UnitPrice     int    `json:"unit_price" db:"unit_price"`
	Quantity      int    `json:"quantity" db:"quantity"`
	Subtotal      int    `json:"subtotal" db:"subtotal"`
}

// CheckoutRequest represents a checkout request payload
// Items are validated in the API and database layers.
type CheckoutRequest struct {
	Items []CheckoutItem `json:"items"`
}

// CheckoutItem represents a product purchase line
// ProductID is required and Quantity must be > 0.
type CheckoutItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

// Checkout creates a transaction, updates product stocks, and inserts transaction details atomically.
func Checkout(db *sql.DB, productTable, categoryTable, transactionTable, transactionDetailTable string, items []CheckoutItem) (Transaction, error) {
	if len(items) == 0 {
		return Transaction{}, ErrCheckoutEmptyItems
	}

	tx, err := db.Begin()
	if err != nil {
		return Transaction{}, fmt.Errorf("failed to begin transaction: %w", err)
	}

	rollback := func() {
		_ = tx.Rollback()
	}

	getProductQuery := fmt.Sprintf("SELECT p.id, p.name, p.price, p.stock, COALESCE(c.description, '') FROM %s p LEFT JOIN %s c ON p.category_id = c.id WHERE p.id = $1 FOR UPDATE OF p", productTable, categoryTable)
	updateStockQuery := fmt.Sprintf("UPDATE %s SET stock = $1 WHERE id = $2", productTable)

	var details []TransactionDetail
	totalAmount := 0

	for _, item := range items {
		if item.ProductID <= 0 || item.Quantity <= 0 {
			rollback()
			return Transaction{}, ErrInvalidCheckoutItem
		}

		var (
			productID   int
			productName string
			productDesc string
			price       int
			stock       int
		)

		err = tx.QueryRow(getProductQuery, item.ProductID).Scan(&productID, &productName, &price, &stock, &productDesc)
		if err != nil {
			rollback()
			if errors.Is(err, sql.ErrNoRows) {
				return Transaction{}, ErrProductNotFound
			}
			return Transaction{}, fmt.Errorf("failed to fetch product: %w", err)
		}

		// Validation: with FOR UPDATE lock, ensure stock is enough to avoid oversell
		if stock < item.Quantity {
			rollback()
			return Transaction{}, ErrInsufficientStock
		}

		newStock := stock - item.Quantity
		_, err = tx.Exec(updateStockQuery, newStock, productID)
		if err != nil {
			rollback()
			return Transaction{}, fmt.Errorf("failed to update stock: %w", err)
		}

		subtotal := price * item.Quantity
		totalAmount += subtotal

		details = append(details, TransactionDetail{
			ProductID:   productID,
			ProductName: productName,
			ProductDesc: productDesc,
			UnitPrice:   price,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	insertTransactionQuery := fmt.Sprintf("INSERT INTO %s (total_amount) VALUES ($1) RETURNING id, total_amount, created_at", transactionTable)
	var transaction Transaction
	transaction.Details = details
	transaction.TotalAmount = totalAmount

	err = tx.QueryRow(insertTransactionQuery, totalAmount).Scan(&transaction.ID, &transaction.TotalAmount, &transaction.CreatedAt)
	if err != nil {
		rollback()
		return Transaction{}, fmt.Errorf("failed to create transaction: %w", err)
	}

	insertDetailQuery := fmt.Sprintf("INSERT INTO %s (transaction_id, product_id, product_name, product_description, unit_price, quantity, subtotal) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id", transactionDetailTable)
	for i := range transaction.Details {
		detail := &transaction.Details[i]
		detail.TransactionID = transaction.ID
		err = tx.QueryRow(insertDetailQuery, transaction.ID, detail.ProductID, detail.ProductName, detail.ProductDesc, detail.UnitPrice, detail.Quantity, detail.Subtotal).Scan(&detail.ID)
		if err != nil {
			rollback()
			return Transaction{}, fmt.Errorf("failed to create transaction detail: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		rollback()
		return Transaction{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return transaction, nil
}
