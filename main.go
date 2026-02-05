package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"codewithumam-tugas1/api"
	"codewithumam-tugas1/config"
	"codewithumam-tugas1/database"
)

// GitCommit will be set at build time
var GitCommit string

func main() {
	// Load configuration from secrets.yml or environment variables
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Build connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.URL, cfg.DBPort, cfg.User, cfg.Password, cfg.Name)

	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Connected to database successfully")

	// Run migrations
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Database migrations completed")

	// Initialize categories service
	categories := api.NewCategories(db, "category")

	// Initialize products service
	products := api.NewProducts(db, "product")

	// Initialize checkout service
	checkout := api.NewCheckout(db, "product", "category", "\"transaction\"", "transaction_detail")

	// Initialize report service
	report := api.NewReport(db, "\"transaction\"", "transaction_detail")

	// Category routes
	http.HandleFunc("GET /categories", categories.GetAll)
	http.HandleFunc("GET /categories/{id}", categories.GetByID)
	http.HandleFunc("POST /categories", categories.Create)
	http.HandleFunc("PUT /categories/{id}", categories.Update)
	http.HandleFunc("DELETE /categories/{id}", categories.Delete)

	// Product routes
	http.HandleFunc("GET /products", products.GetAll)
	http.HandleFunc("GET /products/{id}", products.GetByID)
	http.HandleFunc("POST /products", products.Create)
	http.HandleFunc("PUT /products/{id}", products.Update)
	http.HandleFunc("DELETE /products/{id}", products.Delete)

	// Checkout routes
	http.HandleFunc("POST /checkout", checkout.Create)

	// Report routes
	http.HandleFunc("GET /report/hari-ini", report.Today)
	http.HandleFunc("GET /report", report.Range)

	// Original endpoints
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if GitCommit == "" {
			GitCommit = "unknown"
		}
		fmt.Fprintf(w, "Commit version: %s", GitCommit)
	})

	http.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(addr, nil))
}
