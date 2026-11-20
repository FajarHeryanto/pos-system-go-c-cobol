// ==================================================================
// utils/csv_exporter.go — Export Data Transaksi ke CSV
// Sistem Informasi Penjualan (Point of Sales)
// ==================================================================
//
// Format CSV sesuai Bab 4.3 RPD (Go → C):
//   transaction_id,product_id,product_name,category,quantity,
//   unit_price,discount_percent,payment_method,transaction_date,
//   cashier_name
// ==================================================================

package utils

import (
	"database/sql"
	"fmt"
	"os"
)

// ExportTransactionsToCSV mengekspor semua data transaksi dari SQLite
// ke file CSV sesuai format yang diharapkan oleh program C processor.
//
// Query melakukan JOIN ke tabel products untuk mendapatkan field category
// yang dibutuhkan oleh C processor.
func ExportTransactionsToCSV(db *sql.DB, filepath string) (int, error) {
	// Query: ambil semua transaksi beserta kategori dari tabel products
	rows, err := db.Query(`
		SELECT t.transaction_id, t.product_id, t.product_name,
		       p.category, t.quantity, t.unit_price,
		       t.discount_percent, t.payment_method,
		       t.transaction_date, t.cashier_name
		FROM transactions t
		LEFT JOIN products p ON t.product_id = p.product_id
		ORDER BY t.transaction_date, t.transaction_id
	`)
	if err != nil {
		return 0, fmt.Errorf("gagal query transaksi: %w", err)
	}
	defer rows.Close()

	// Buka file output
	file, err := os.Create(filepath)
	if err != nil {
		return 0, fmt.Errorf("gagal membuat file %s: %w", filepath, err)
	}
	defer file.Close()

	// Tulis header CSV
	_, err = fmt.Fprintln(file, "transaction_id,product_id,product_name,category,quantity,unit_price,discount_percent,payment_method,transaction_date,cashier_name")
	if err != nil {
		return 0, fmt.Errorf("gagal menulis header CSV: %w", err)
	}

	// Tulis setiap baris data
	count := 0
	for rows.Next() {
		var (
			transID, prodID, prodName, payMethod, transDate, cashier string
			category                                                 sql.NullString
			qty                                                      int
			unitPrice, discountPct                                   float64
		)

		err := rows.Scan(&transID, &prodID, &prodName, &category,
			&qty, &unitPrice, &discountPct, &payMethod, &transDate, &cashier)
		if err != nil {
			return count, fmt.Errorf("gagal scan baris ke-%d: %w", count+1, err)
		}

		cat := "Unknown"
		if category.Valid {
			cat = category.String
		}

		_, err = fmt.Fprintf(file, "%s,%s,%s,%s,%d,%.0f,%.0f,%s,%s,%s\n",
			transID, prodID, prodName, cat, qty, unitPrice, discountPct,
			payMethod, transDate, cashier)
		if err != nil {
			return count, fmt.Errorf("gagal menulis baris ke-%d: %w", count+1, err)
		}
		count++
	}

	if err := rows.Err(); err != nil {
		return count, fmt.Errorf("error iterasi rows: %w", err)
	}

	return count, nil
}
