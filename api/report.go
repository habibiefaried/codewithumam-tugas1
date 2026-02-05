package api

import (
	"database/sql"
	"net/http"
	"time"

	"codewithumam-tugas1/database"
	"encoding/json"
)

// Report handles report endpoints.
type Report struct {
	db                     *sql.DB
	transactionTable       string
	transactionDetailTable string
}

// NewReport creates a new report service.
func NewReport(db *sql.DB, transactionTable, transactionDetailTable string) *Report {
	return &Report{
		db:                     db,
		transactionTable:       transactionTable,
		transactionDetailTable: transactionDetailTable,
	}
}

// Today handles GET /report/hari-ini
func (r *Report) Today(w http.ResponseWriter, _ *http.Request) {
	summary, err := database.GetReportToday(r.db, r.transactionTable, r.transactionDetailTable, time.Now().UTC())
	if err != nil {
		http.Error(w, "Failed to generate report", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

// Range handles GET /report?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD
func (r *Report) Range(w http.ResponseWriter, req *http.Request) {
	startDate := req.URL.Query().Get("start_date")
	endDate := req.URL.Query().Get("end_date")
	if startDate == "" || endDate == "" {
		http.Error(w, "start_date and end_date are required", http.StatusBadRequest)
		return
	}

	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		http.Error(w, "Invalid start_date", http.StatusBadRequest)
		return
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		http.Error(w, "Invalid end_date", http.StatusBadRequest)
		return
	}
	if end.Before(start) {
		http.Error(w, "end_date must be on or after start_date", http.StatusBadRequest)
		return
	}

	endExclusive := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, 1)
	startUTC := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)

	summary, err := database.GetReportBetween(r.db, r.transactionTable, r.transactionDetailTable, startUTC, endExclusive)
	if err != nil {
		http.Error(w, "Failed to generate report", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}
