package database

import "testing"

func TestCheckoutSuccess(t *testing.T) {
	db := setupProductTestDB(t)
	defer teardownProductTestDB(t, db)

	cat, err := Create(db, "category_test", "Food", "Food category")
	if err != nil {
		t.Fatalf("Failed to create category: %v", err)
	}

	prod1, err := CreateProduct(db, "product_test", "category_test", "Apple", 10, 50, cat.ID)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	prod2, err := CreateProduct(db, "product_test", "category_test", "Orange", 20, 30, cat.ID)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	items := []CheckoutItem{
		{ProductID: prod1.ID, Quantity: 3},
		{ProductID: prod2.ID, Quantity: 2},
	}

	trx, err := Checkout(db, "product_test", "category_test", "transaction_test", "transaction_detail_test", items)
	if err != nil {
		t.Fatalf("Checkout failed: %v", err)
	}

	if trx.ID == 0 {
		t.Error("Expected non-zero transaction ID")
	}

	if trx.TotalAmount != (3*10 + 2*20) {
		t.Errorf("Expected total amount 70, got %d", trx.TotalAmount)
	}

	if trx.CreatedAt.IsZero() {
		t.Error("Expected created_at to be set")
	}

	if len(trx.Details) != 2 {
		t.Fatalf("Expected 2 details, got %d", len(trx.Details))
	}
	if trx.Details[0].UnitPrice != 10 {
		t.Errorf("Expected unit price 10, got %d", trx.Details[0].UnitPrice)
	}
	if trx.Details[0].ProductDesc != "Food category" {
		t.Errorf("Expected product description 'Food category', got %s", trx.Details[0].ProductDesc)
	}

	updated1, err := GetProductByID(db, "product_test", "category_test", prod1.ID)
	if err != nil {
		t.Fatalf("Failed to fetch product after checkout: %v", err)
	}
	if updated1.Stock != 47 {
		t.Errorf("Expected stock 47, got %d", updated1.Stock)
	}

	updated2, err := GetProductByID(db, "product_test", "category_test", prod2.ID)
	if err != nil {
		t.Fatalf("Failed to fetch product after checkout: %v", err)
	}
	if updated2.Stock != 28 {
		t.Errorf("Expected stock 28, got %d", updated2.Stock)
	}
}

func TestCheckoutWithoutCategory(t *testing.T) {
	db := setupProductTestDB(t)
	defer teardownProductTestDB(t, db)

	prod, err := CreateProduct(db, "product_test", "category_test", "NoCatItem", 100, 5, 0)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	trx, err := Checkout(db, "product_test", "category_test", "transaction_test", "transaction_detail_test", []CheckoutItem{{ProductID: prod.ID, Quantity: 1}})
	if err != nil {
		t.Fatalf("Checkout failed: %v", err)
	}
	if len(trx.Details) != 1 {
		t.Fatalf("Expected 1 detail, got %d", len(trx.Details))
	}
	if trx.Details[0].ProductDesc != "" {
		t.Errorf("Expected empty product description, got %s", trx.Details[0].ProductDesc)
	}
}

func TestCheckoutInsufficientStock(t *testing.T) {
	db := setupProductTestDB(t)
	defer teardownProductTestDB(t, db)

	cat, err := Create(db, "category_test", "Food", "Food category")
	if err != nil {
		t.Fatalf("Failed to create category: %v", err)
	}

	prod, err := CreateProduct(db, "product_test", "category_test", "Apple", 10, 1, cat.ID)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	_, err = Checkout(db, "product_test", "category_test", "transaction_test", "transaction_detail_test", []CheckoutItem{{ProductID: prod.ID, Quantity: 2}})
	if err == nil {
		t.Fatal("Expected insufficient stock error")
	}
	if err != ErrInsufficientStock {
		t.Fatalf("Expected ErrInsufficientStock, got %v", err)
	}
}

func TestCheckoutProductNotFound(t *testing.T) {
	db := setupProductTestDB(t)
	defer teardownProductTestDB(t, db)

	_, err := Checkout(db, "product_test", "category_test", "transaction_test", "transaction_detail_test", []CheckoutItem{{ProductID: 9999, Quantity: 1}})
	if err == nil {
		t.Fatal("Expected product not found error")
	}
	if err != ErrProductNotFound {
		t.Fatalf("Expected ErrProductNotFound, got %v", err)
	}
}

func TestCheckoutInvalidItems(t *testing.T) {
	db := setupProductTestDB(t)
	defer teardownProductTestDB(t, db)

	_, err := Checkout(db, "product_test", "category_test", "transaction_test", "transaction_detail_test", []CheckoutItem{})
	if err == nil {
		t.Fatal("Expected empty items error")
	}
	if err != ErrCheckoutEmptyItems {
		t.Fatalf("Expected ErrCheckoutEmptyItems, got %v", err)
	}

	_, err = Checkout(db, "product_test", "category_test", "transaction_test", "transaction_detail_test", []CheckoutItem{{ProductID: 1, Quantity: 0}})
	if err == nil {
		t.Fatal("Expected invalid item error")
	}
	if err != ErrInvalidCheckoutItem {
		t.Fatalf("Expected ErrInvalidCheckoutItem, got %v", err)
	}

	_, err = Checkout(db, "product_test", "category_test", "transaction_test", "transaction_detail_test", []CheckoutItem{{ProductID: 0, Quantity: 1}})
	if err == nil {
		t.Fatal("Expected invalid item error for product_id 0")
	}
	if err != ErrInvalidCheckoutItem {
		t.Fatalf("Expected ErrInvalidCheckoutItem, got %v", err)
	}

	_, err = Checkout(db, "product_test", "category_test", "transaction_test", "transaction_detail_test", []CheckoutItem{{ProductID: 1, Quantity: -2}})
	if err == nil {
		t.Fatal("Expected invalid item error for negative quantity")
	}
	if err != ErrInvalidCheckoutItem {
		t.Fatalf("Expected ErrInvalidCheckoutItem, got %v", err)
	}
}

func TestCheckoutRollbackOnInvalidItem(t *testing.T) {
	db := setupProductTestDB(t)
	defer teardownProductTestDB(t, db)

	cat, err := Create(db, "category_test", "Food", "Food category")
	if err != nil {
		t.Fatalf("Failed to create category: %v", err)
	}

	prod, err := CreateProduct(db, "product_test", "category_test", "Apple", 10, 10, cat.ID)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	_, err = Checkout(db, "product_test", "category_test", "transaction_test", "transaction_detail_test", []CheckoutItem{
		{ProductID: prod.ID, Quantity: 2},
		{ProductID: 0, Quantity: 1},
	})
	if err == nil {
		t.Fatal("Expected invalid item error")
	}
	if err != ErrInvalidCheckoutItem {
		t.Fatalf("Expected ErrInvalidCheckoutItem, got %v", err)
	}

	updated, err := GetProductByID(db, "product_test", "category_test", prod.ID)
	if err != nil {
		t.Fatalf("Failed to fetch product after rollback: %v", err)
	}
	if updated.Stock != 10 {
		t.Errorf("Expected stock 10 after rollback, got %d", updated.Stock)
	}
}

func TestCheckoutRollbackOnInsufficientStockSecondItem(t *testing.T) {
	db := setupProductTestDB(t)
	defer teardownProductTestDB(t, db)

	cat, err := Create(db, "category_test", "Food", "Food category")
	if err != nil {
		t.Fatalf("Failed to create category: %v", err)
	}

	prod1, err := CreateProduct(db, "product_test", "category_test", "Apple", 10, 10, cat.ID)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	prod2, err := CreateProduct(db, "product_test", "category_test", "Orange", 20, 1, cat.ID)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	_, err = Checkout(db, "product_test", "category_test", "transaction_test", "transaction_detail_test", []CheckoutItem{
		{ProductID: prod1.ID, Quantity: 2},
		{ProductID: prod2.ID, Quantity: 2},
	})
	if err == nil {
		t.Fatal("Expected insufficient stock error")
	}
	if err != ErrInsufficientStock {
		t.Fatalf("Expected ErrInsufficientStock, got %v", err)
	}

	updated1, err := GetProductByID(db, "product_test", "category_test", prod1.ID)
	if err != nil {
		t.Fatalf("Failed to fetch product after rollback: %v", err)
	}
	if updated1.Stock != 10 {
		t.Errorf("Expected stock 10 after rollback, got %d", updated1.Stock)
	}

	updated2, err := GetProductByID(db, "product_test", "category_test", prod2.ID)
	if err != nil {
		t.Fatalf("Failed to fetch product after rollback: %v", err)
	}
	if updated2.Stock != 1 {
		t.Errorf("Expected stock 1 after rollback, got %d", updated2.Stock)
	}
}

func TestCheckoutRollbackOnProductNotFoundSecondItem(t *testing.T) {
	db := setupProductTestDB(t)
	defer teardownProductTestDB(t, db)

	cat, err := Create(db, "category_test", "Food", "Food category")
	if err != nil {
		t.Fatalf("Failed to create category: %v", err)
	}

	prod, err := CreateProduct(db, "product_test", "category_test", "Apple", 10, 10, cat.ID)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	_, err = Checkout(db, "product_test", "category_test", "transaction_test", "transaction_detail_test", []CheckoutItem{
		{ProductID: prod.ID, Quantity: 2},
		{ProductID: 9999, Quantity: 1},
	})
	if err == nil {
		t.Fatal("Expected product not found error")
	}
	if err != ErrProductNotFound {
		t.Fatalf("Expected ErrProductNotFound, got %v", err)
	}

	updated, err := GetProductByID(db, "product_test", "category_test", prod.ID)
	if err != nil {
		t.Fatalf("Failed to fetch product after rollback: %v", err)
	}
	if updated.Stock != 10 {
		t.Errorf("Expected stock 10 after rollback, got %d", updated.Stock)
	}
}
