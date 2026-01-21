package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// This will be set at build time
var GitCommit string

func main() {
	// Get port from environment variable or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Handler function
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!")
	})

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})

	// Version endpoint
	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		if GitCommit == "" {
			GitCommit = "unknown"
		}
		fmt.Fprintf(w, "Commit: %s", GitCommit)
	})

	// Start server
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(addr, nil))
}
