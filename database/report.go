package database

import (
	"database/sql"
	"fmt"
	"time"
)

// ReportTopProduct represents the best-selling product
// based on total quantity sold.
type ReportTopProduct struct {
	Nama       string `json:"nama"`
	QtyTerjual int    `json:"qty_terjual"`
}

// ReportSummary represents revenue and transaction aggregates.
type ReportSummary struct {
	TotalRevenue   int              `json:"total_revenue"`
	TotalTransaksi int              `json:"total_transaksi"`
	ProdukTerlaris ReportTopProduct `json:"produk_terlaris"`
}

// GetReportBetween aggregates revenue, transaction count, and top product within a date range.
// Range is [start, end), so pass end as the next day for inclusive end-date.
func GetReportBetween(db *sql.DB, transactionTable, transactionDetailTable string, start, end time.Time) (ReportSummary, error) {
	summary := ReportSummary{}

	aggQuery := fmt.Sprintf("SELECT COALESCE(SUM(total_amount), 0), COUNT(*) FROM %s WHERE created_at >= $1 AND created_at < $2", transactionTable)
	err := db.QueryRow(aggQuery, start, end).Scan(&summary.TotalRevenue, &summary.TotalTransaksi)
	if err != nil {
		return ReportSummary{}, fmt.Errorf("failed to aggregate transactions: %w", err)
	}

	topQuery := fmt.Sprintf("SELECT product_name, COALESCE(SUM(quantity),0) FROM %s d JOIN %s t ON d.transaction_id = t.id WHERE t.created_at >= $1 AND t.created_at < $2 GROUP BY product_name ORDER BY SUM(quantity) DESC LIMIT 1", transactionDetailTable, transactionTable)
	var topName sql.NullString
	var topQty sql.NullInt64
	err = db.QueryRow(topQuery, start, end).Scan(&topName, &topQty)
	if err != nil && err != sql.ErrNoRows {
		return ReportSummary{}, fmt.Errorf("failed to aggregate top product: %w", err)
	}
	if topName.Valid {
		summary.ProdukTerlaris.Nama = topName.String
		summary.ProdukTerlaris.QtyTerjual = int(topQty.Int64)
	}

	return summary, nil
}

// GetReportToday aggregates report for the given day based on UTC date boundaries.
func GetReportToday(db *sql.DB, transactionTable, transactionDetailTable string, day time.Time) (ReportSummary, error) {
	start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 0, 1)
	return GetReportBetween(db, transactionTable, transactionDetailTable, start, end)
}
