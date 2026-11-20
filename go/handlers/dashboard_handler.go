// ==================================================================
// handlers/dashboard_handler.go — Handler Data Dashboard
// Sistem Informasi Penjualan (Point of Sales)
// ==================================================================
//
// Endpoint: GET /api/dashboard
//
// Mengembalikan ringkasan statistik:
//   - Total produk
//   - Total transaksi
//   - Total pendapatan (grand total)
//   - Jumlah kasir aktif
//   - Distribusi per kategori
//   - 5 transaksi terakhir
//   - 5 produk terlaris (by total qty)
// ==================================================================

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"pos-system/database"
	"pos-system/utils"
)

// DashboardStats merepresentasikan data statistik untuk dashboard.
type DashboardStats struct {
	TotalProducts     int                `json:"total_products"`
	TotalTransactions int                `json:"total_transactions"`
	TotalRevenue      float64            `json:"total_revenue"`
	TotalCashiers     int                `json:"total_cashiers"`
	CategoryStats     []CategoryStatItem `json:"category_stats"`
	RecentTrans       []RecentTransItem  `json:"recent_transactions"`
	TopProducts       []TopProductItem   `json:"top_products"`
}

// CategoryStatItem merepresentasikan statistik per kategori.
type CategoryStatItem struct {
	Category string  `json:"category"`
	Count    int     `json:"count"`
	Revenue  float64 `json:"revenue"`
}

// RecentTransItem merepresentasikan satu transaksi terbaru.
type RecentTransItem struct {
	TransactionID string  `json:"transaction_id"`
	ProductName   string  `json:"product_name"`
	Quantity      int     `json:"quantity"`
	TotalPrice    float64 `json:"total_price"`
	Date          string  `json:"transaction_date"`
}

// TopProductItem merepresentasikan produk terlaris.
type TopProductItem struct {
	ProductName string  `json:"product_name"`
	TotalQty    int     `json:"total_qty"`
	TotalRev    float64 `json:"total_revenue"`
}

// GetDashboardStats mengambil data ringkasan untuk dashboard.
func GetDashboardStats(c *gin.Context) {
	var stats DashboardStats

	// Total produk
	database.DB.QueryRow("SELECT COUNT(*) FROM products").Scan(&stats.TotalProducts)

	// Total transaksi
	database.DB.QueryRow("SELECT COUNT(*) FROM transactions").Scan(&stats.TotalTransactions)

	// Total pendapatan
	database.DB.QueryRow("SELECT COALESCE(SUM(total_price), 0) FROM transactions").Scan(&stats.TotalRevenue)

	// Jumlah kasir unik
	database.DB.QueryRow("SELECT COUNT(DISTINCT cashier_name) FROM transactions").Scan(&stats.TotalCashiers)

	// Statistik per kategori
	catRows, err := database.DB.Query(`
		SELECT p.category, COUNT(*) as cnt, COALESCE(SUM(t.total_price), 0) as rev
		FROM transactions t
		JOIN products p ON t.product_id = p.product_id
		GROUP BY p.category
		ORDER BY rev DESC
	`)
	if err == nil {
		defer catRows.Close()
		for catRows.Next() {
			var item CategoryStatItem
			if catRows.Scan(&item.Category, &item.Count, &item.Revenue) == nil {
				stats.CategoryStats = append(stats.CategoryStats, item)
			}
		}
	}
	if stats.CategoryStats == nil {
		stats.CategoryStats = []CategoryStatItem{}
	}

	// 5 transaksi terakhir
	recentRows, err := database.DB.Query(`
		SELECT transaction_id, product_name, quantity, total_price, transaction_date
		FROM transactions
		ORDER BY transaction_date DESC, transaction_id DESC
		LIMIT 5
	`)
	if err == nil {
		defer recentRows.Close()
		for recentRows.Next() {
			var item RecentTransItem
			if recentRows.Scan(&item.TransactionID, &item.ProductName,
				&item.Quantity, &item.TotalPrice, &item.Date) == nil {
				stats.RecentTrans = append(stats.RecentTrans, item)
			}
		}
	}
	if stats.RecentTrans == nil {
		stats.RecentTrans = []RecentTransItem{}
	}

	// Top 5 produk (berdasarkan total quantity terjual)
	topRows, err := database.DB.Query(`
		SELECT product_name, SUM(quantity) as total_qty, SUM(total_price) as total_rev
		FROM transactions
		GROUP BY product_name
		ORDER BY total_qty DESC
		LIMIT 5
	`)
	if err == nil {
		defer topRows.Close()
		for topRows.Next() {
			var item TopProductItem
			if topRows.Scan(&item.ProductName, &item.TotalQty, &item.TotalRev) == nil {
				stats.TopProducts = append(stats.TopProducts, item)
			}
		}
	}
	if stats.TopProducts == nil {
		stats.TopProducts = []TopProductItem{}
	}

	utils.RespondSuccess(c, http.StatusOK, "", stats)
}
