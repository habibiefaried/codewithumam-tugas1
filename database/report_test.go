package database

import (
	"database/sql"
	"testing"
	"time"
)

func TestGetReportBetween(t *testing.T) {
	db := setupProductTestDB(t)
	defer teardownProductTestDB(t, db)

	cat, err := Create(db, "category_test", "Food", "Food category")
	if err != nil {
		t.Fatalf("Failed to create category: %v", err)
	}

	prod1, err := CreateProduct(db, "product_test", "category_test", "Indomie Goreng", 5000, 100, cat.ID)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}
	prod2, err := CreateProduct(db, "product_test", "category_test", "Teh Botol", 3000, 100, cat.ID)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	rangeStart := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	rangeEnd := time.Date(2026, 2, 2, 0, 0, 0, 0, time.UTC)

	trx1 := insertTransactionWithDetails(t, db, 15000, rangeStart.Add(2*time.Hour), []TransactionDetail{
		{ProductID: prod1.ID, ProductName: prod1.Name, Quantity: 2, Subtotal: 10000},
		{ProductID: prod2.ID, ProductName: prod2.Name, Quantity: 1, Subtotal: 5000},
	})
	_ = trx1

	trx2 := insertTransactionWithDetails(t, db, 20000, rangeStart.Add(5*time.Hour), []TransactionDetail{
		{ProductID: prod1.ID, ProductName: prod1.Name, Quantity: 3, Subtotal: 15000},
		{ProductID: prod2.ID, ProductName: prod2.Name, Quantity: 2, Subtotal: 5000},
	})
	_ = trx2

	// Outside range
	_ = insertTransactionWithDetails(t, db, 9999, rangeStart.AddDate(0, 0, -1), []TransactionDetail{
		{ProductID: prod1.ID, ProductName: prod1.Name, Quantity: 1, Subtotal: 9999},
	})

	summary, err := GetReportBetween(db, "transaction_test", "transaction_detail_test", rangeStart, rangeEnd)
	if err != nil {
		t.Fatalf("GetReportBetween failed: %v", err)
	}

	if summary.TotalRevenue != 35000 {
		t.Errorf("Expected total_revenue 35000, got %d", summary.TotalRevenue)
	}
	if summary.TotalTransaksi != 2 {
		t.Errorf("Expected total_transaksi 2, got %d", summary.TotalTransaksi)
	}
	if summary.ProdukTerlaris.Nama != "Indomie Goreng" {
		t.Errorf("Expected produk_terlaris Indomie Goreng, got %s", summary.ProdukTerlaris.Nama)
	}
	if summary.ProdukTerlaris.QtyTerjual != 5 {
		t.Errorf("Expected qty_terjual 5, got %d", summary.ProdukTerlaris.QtyTerjual)
	}
}

func TestGetReportBetweenEmpty(t *testing.T) {
	db := setupProductTestDB(t)
	defer teardownProductTestDB(t, db)

	rangeStart := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	rangeEnd := time.Date(2026, 2, 2, 0, 0, 0, 0, time.UTC)

	summary, err := GetReportBetween(db, "transaction_test", "transaction_detail_test", rangeStart, rangeEnd)
	if err != nil {
		t.Fatalf("GetReportBetween failed: %v", err)
	}

	if summary.TotalRevenue != 0 {
		t.Errorf("Expected total_revenue 0, got %d", summary.TotalRevenue)
	}
	if summary.TotalTransaksi != 0 {
		t.Errorf("Expected total_transaksi 0, got %d", summary.TotalTransaksi)
	}
	if summary.ProdukTerlaris.Nama != "" {
		t.Errorf("Expected empty produk_terlaris name, got %s", summary.ProdukTerlaris.Nama)
	}
	if summary.ProdukTerlaris.QtyTerjual != 0 {
		t.Errorf("Expected qty_terjual 0, got %d", summary.ProdukTerlaris.QtyTerjual)
	}
}

func TestGetReportToday(t *testing.T) {
	db := setupProductTestDB(t)
	defer teardownProductTestDB(t, db)

	cat, err := Create(db, "category_test", "Food", "Food category")
	if err != nil {
		t.Fatalf("Failed to create category: %v", err)
	}

	prod, err := CreateProduct(db, "product_test", "category_test", "Indomie Goreng", 5000, 100, cat.ID)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	day := time.Date(2026, 2, 5, 10, 0, 0, 0, time.UTC)
	_ = insertTransactionWithDetails(t, db, 10000, day.Add(2*time.Hour), []TransactionDetail{
		{ProductID: prod.ID, ProductName: prod.Name, Quantity: 2, Subtotal: 10000},
	})

	summary, err := GetReportToday(db, "transaction_test", "transaction_detail_test", day)
	if err != nil {
		t.Fatalf("GetReportToday failed: %v", err)
	}

	if summary.TotalRevenue != 10000 {
		t.Errorf("Expected total_revenue 10000, got %d", summary.TotalRevenue)
	}
	if summary.TotalTransaksi != 1 {
		t.Errorf("Expected total_transaksi 1, got %d", summary.TotalTransaksi)
	}
}

func insertTransactionWithDetails(t *testing.T, db *sql.DB, total int, createdAt time.Time, details []TransactionDetail) int {
	var trxID int
	err := db.QueryRow("INSERT INTO transaction_test (total_amount, created_at) VALUES ($1, $2) RETURNING id", total, createdAt).Scan(&trxID)
	if err != nil {
		t.Fatalf("Failed to insert transaction_test: %v", err)
	}

	for _, d := range details {
		_, err = db.Exec("INSERT INTO transaction_detail_test (transaction_id, product_id, product_name, product_description, unit_price, quantity, subtotal) VALUES ($1, $2, $3, $4, $5, $6, $7)", trxID, d.ProductID, d.ProductName, d.ProductDesc, d.UnitPrice, d.Quantity, d.Subtotal)
		if err != nil {
			t.Fatalf("Failed to insert transaction_detail_test: %v", err)
		}
	}
	return trxID
}
