package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"codewithumam-tugas1/config"
	"codewithumam-tugas1/database"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.URL, cfg.DBPort, cfg.User, cfg.Password, cfg.Name)

	log.Printf("Connecting to database host=%s port=%s dbname=%s user=%s", cfg.URL, cfg.DBPort, cfg.Name, cfg.User)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	dropSQL := `
	DROP TABLE IF EXISTS transaction_detail_test;
	DROP TABLE IF EXISTS transaction_test;
	DROP TABLE IF EXISTS product_test;
	DROP TABLE IF EXISTS category_test;
	DROP TABLE IF EXISTS transaction_detail;
	DROP TABLE IF EXISTS "transaction";
	DROP TABLE IF EXISTS product;
	DROP TABLE IF EXISTS category;
	`

	if _, err := db.Exec(dropSQL); err != nil {
		log.Fatalf("Failed to drop tables: %v", err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Database reset completed")
}
