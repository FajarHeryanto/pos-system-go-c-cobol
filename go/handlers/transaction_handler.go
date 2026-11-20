// ==================================================================
// handlers/transaction_handler.go — CRUD Transaksi Handler
// Sistem Informasi Penjualan (Point of Sales)
// ==================================================================
//
// Endpoint:
//   GET    /api/transactions      → GetAllTransactions
//   POST   /api/transactions      → CreateTransaction
//   DELETE /api/transactions/:id  → DeleteTransaction
// ==================================================================

package handlers

import (
	"database/sql"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"pos-system/database"
	"pos-system/models"
	"pos-system/utils"
)

// GetAllTransactions mengembalikan seluruh daftar transaksi.
func GetAllTransactions(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT transaction_id, product_id, product_name, quantity,
		       unit_price, discount_percent, total_price,
		       payment_method, transaction_date, cashier_name
		FROM transactions
		ORDER BY transaction_date DESC, transaction_id DESC
	`)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Gagal mengambil data transaksi: "+err.Error())
		return
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.TransactionID, &t.ProductID, &t.ProductName,
			&t.Quantity, &t.UnitPrice, &t.DiscountPercent, &t.TotalPrice,
			&t.PaymentMethod, &t.TransactionDate, &t.CashierName); err != nil {
			utils.RespondError(c, http.StatusInternalServerError, "Gagal membaca data transaksi: "+err.Error())
			return
		}
		transactions = append(transactions, t)
	}

	if transactions == nil {
		transactions = []models.Transaction{}
	}

	utils.RespondSuccess(c, http.StatusOK, "", transactions)
}

// CreateTransaction membuat transaksi baru.
// Proses:
//  1. Validasi input
//  2. Cek produk ada dan stok mencukupi
//  3. Hitung total_price (qty * unit_price * (1 - discount/100))
//  4. Insert transaksi dan kurangi stok produk
func CreateTransaction(c *gin.Context) {
	var input models.TransactionInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Format request tidak valid")
		return
	}

	// Validasi input
	if validationErrors := utils.ValidateTransactionInput(input); len(validationErrors) > 0 {
		utils.RespondValidationError(c, validationErrors)
		return
	}

	// Cek produk dan stok
	var product models.Product
	err := database.DB.QueryRow(
		"SELECT product_id, name, category, price, stock FROM products WHERE product_id = ?",
		input.ProductID,
	).Scan(&product.ProductID, &product.Name, &product.Category, &product.Price, &product.Stock)

	if err == sql.ErrNoRows {
		utils.RespondError(c, http.StatusNotFound, "Produk tidak ditemukan")
		return
	}
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Gagal mengambil data produk: "+err.Error())
		return
	}

	// Cek stok mencukupi
	if product.Stock < input.Quantity {
		utils.RespondError(c, http.StatusBadRequest,
			fmt.Sprintf("Stok tidak mencukupi. Stok tersedia: %d, diminta: %d", product.Stock, input.Quantity))
		return
	}

	// Generate transaction_id berikutnya
	var maxID sql.NullString
	database.DB.QueryRow("SELECT MAX(transaction_id) FROM transactions").Scan(&maxID)

	newID := "T001"
	if maxID.Valid && len(maxID.String) >= 2 {
		var num int
		fmt.Sscanf(maxID.String[1:], "%d", &num)
		newID = fmt.Sprintf("T%03d", num+1)
	}

	// Hitung total_price
	subtotal := float64(input.Quantity) * product.Price
	discountAmount := subtotal * input.DiscountPercent / 100.0
	totalPrice := math.Round(subtotal - discountAmount)

	// Tanggal transaksi (hari ini)
	transDate := time.Now().Format("2006-01-02")

	// Begin transaction (database)
	tx, err := database.DB.Begin()
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Gagal memulai transaksi database")
		return
	}

	// Insert transaksi
	_, err = tx.Exec(`
		INSERT INTO transactions 
		(transaction_id, product_id, product_name, quantity, unit_price,
		 discount_percent, total_price, payment_method, transaction_date, cashier_name)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		newID, input.ProductID, product.Name, input.Quantity, product.Price,
		input.DiscountPercent, totalPrice, input.PaymentMethod, transDate, input.CashierName,
	)
	if err != nil {
		tx.Rollback()
		utils.RespondError(c, http.StatusInternalServerError, "Gagal menyimpan transaksi: "+err.Error())
		return
	}

	// Kurangi stok produk
	_, err = tx.Exec("UPDATE products SET stock = stock - ? WHERE product_id = ?",
		input.Quantity, input.ProductID)
	if err != nil {
		tx.Rollback()
		utils.RespondError(c, http.StatusInternalServerError, "Gagal mengupdate stok: "+err.Error())
		return
	}

	// Commit
	if err := tx.Commit(); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Gagal commit transaksi: "+err.Error())
		return
	}

	// Response
	transaction := models.Transaction{
		TransactionID:   newID,
		ProductID:       input.ProductID,
		ProductName:     product.Name,
		Quantity:        input.Quantity,
		UnitPrice:       product.Price,
		DiscountPercent: input.DiscountPercent,
		TotalPrice:      totalPrice,
		PaymentMethod:   input.PaymentMethod,
		TransactionDate: transDate,
		CashierName:     input.CashierName,
	}

	utils.RespondSuccess(c, http.StatusCreated, "Transaksi berhasil disimpan", transaction)
}

// DeleteTransaction menghapus transaksi berdasarkan ID.
// Stok produk dikembalikan (ditambah kembali).
func DeleteTransaction(c *gin.Context) {
	id := c.Param("id")

	// Ambil data transaksi yang akan dihapus (untuk mengembalikan stok)
	var productID string
	var qty int
	err := database.DB.QueryRow(
		"SELECT product_id, quantity FROM transactions WHERE transaction_id = ?", id,
	).Scan(&productID, &qty)

	if err == sql.ErrNoRows {
		utils.RespondError(c, http.StatusNotFound, "Transaksi tidak ditemukan")
		return
	}
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Gagal mengambil data transaksi: "+err.Error())
		return
	}

	// Begin transaction
	tx, err := database.DB.Begin()
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Gagal memulai transaksi database")
		return
	}

	// Hapus transaksi
	_, err = tx.Exec("DELETE FROM transactions WHERE transaction_id = ?", id)
	if err != nil {
		tx.Rollback()
		utils.RespondError(c, http.StatusInternalServerError, "Gagal menghapus transaksi: "+err.Error())
		return
	}

	// Kembalikan stok
	_, err = tx.Exec("UPDATE products SET stock = stock + ? WHERE product_id = ?", qty, productID)
	if err != nil {
		tx.Rollback()
		utils.RespondError(c, http.StatusInternalServerError, "Gagal mengembalikan stok: "+err.Error())
		return
	}

	if err := tx.Commit(); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Gagal commit: "+err.Error())
		return
	}

	utils.RespondSuccess(c, http.StatusOK, "Transaksi berhasil dihapus", nil)
}
