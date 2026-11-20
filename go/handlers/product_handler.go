// ==================================================================
// handlers/product_handler.go — CRUD Produk Handler
// Sistem Informasi Penjualan (Point of Sales)
// ==================================================================
//
// Endpoint:
//   GET    /api/products      → GetAllProducts
//   GET    /api/products/:id  → GetProduct
//   POST   /api/products      → CreateProduct
//   PUT    /api/products/:id  → UpdateProduct
//   DELETE /api/products/:id  → DeleteProduct
// ==================================================================

package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"pos-system/database"
	"pos-system/models"
	"pos-system/utils"
)

// GetAllProducts mengembalikan seluruh daftar produk.
func GetAllProducts(c *gin.Context) {
	rows, err := database.DB.Query(
		"SELECT product_id, name, category, price, stock, created_at FROM products ORDER BY product_id")
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Gagal mengambil data produk: "+err.Error())
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ProductID, &p.Name, &p.Category, &p.Price, &p.Stock, &p.CreatedAt); err != nil {
			utils.RespondError(c, http.StatusInternalServerError, "Gagal membaca data produk: "+err.Error())
			return
		}
		products = append(products, p)
	}

	if products == nil {
		products = []models.Product{}
	}

	utils.RespondSuccess(c, http.StatusOK, "", products)
}

// GetProduct mengembalikan satu produk berdasarkan ID.
func GetProduct(c *gin.Context) {
	id := c.Param("id")

	var p models.Product
	err := database.DB.QueryRow(
		"SELECT product_id, name, category, price, stock, created_at FROM products WHERE product_id = ?", id,
	).Scan(&p.ProductID, &p.Name, &p.Category, &p.Price, &p.Stock, &p.CreatedAt)

	if err == sql.ErrNoRows {
		utils.RespondError(c, http.StatusNotFound, "Produk tidak ditemukan")
		return
	}
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Gagal mengambil data produk: "+err.Error())
		return
	}

	utils.RespondSuccess(c, http.StatusOK, "", p)
}

// CreateProduct menambahkan produk baru ke database.
func CreateProduct(c *gin.Context) {
	var input models.ProductInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Format request tidak valid")
		return
	}

	// Validasi input
	if validationErrors := utils.ValidateProductInput(input); len(validationErrors) > 0 {
		utils.RespondValidationError(c, validationErrors)
		return
	}

	// Generate product_id berikutnya
	var maxID sql.NullString
	database.DB.QueryRow("SELECT MAX(product_id) FROM products").Scan(&maxID)

	newID := "P001"
	if maxID.Valid && len(maxID.String) >= 2 {
		var num int
		fmt.Sscanf(maxID.String[1:], "%d", &num)
		newID = fmt.Sprintf("P%03d", num+1)
	}

	// Insert produk baru
	_, err := database.DB.Exec(
		"INSERT INTO products (product_id, name, category, price, stock) VALUES (?, ?, ?, ?, ?)",
		newID, input.Name, input.Category, input.Price, input.Stock,
	)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Gagal menyimpan produk: "+err.Error())
		return
	}

	// Ambil data yang baru disimpan (termasuk created_at)
	var product models.Product
	database.DB.QueryRow(
		"SELECT product_id, name, category, price, stock, created_at FROM products WHERE product_id = ?", newID,
	).Scan(&product.ProductID, &product.Name, &product.Category, &product.Price, &product.Stock, &product.CreatedAt)

	utils.RespondSuccess(c, http.StatusCreated, "Produk berhasil ditambahkan", product)
}

// UpdateProduct mengubah data produk yang sudah ada.
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	// Cek apakah produk ada
	var exists int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM products WHERE product_id = ?", id).Scan(&exists)
	if err != nil || exists == 0 {
		utils.RespondError(c, http.StatusNotFound, "Produk tidak ditemukan")
		return
	}

	var input models.ProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Format request tidak valid")
		return
	}

	// Validasi input
	if validationErrors := utils.ValidateProductInput(input); len(validationErrors) > 0 {
		utils.RespondValidationError(c, validationErrors)
		return
	}

	// Update produk
	_, err = database.DB.Exec(
		"UPDATE products SET name = ?, category = ?, price = ?, stock = ? WHERE product_id = ?",
		input.Name, input.Category, input.Price, input.Stock, id,
	)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Gagal mengupdate produk: "+err.Error())
		return
	}

	// Ambil data terbaru
	var product models.Product
	database.DB.QueryRow(
		"SELECT product_id, name, category, price, stock, created_at FROM products WHERE product_id = ?", id,
	).Scan(&product.ProductID, &product.Name, &product.Category, &product.Price, &product.Stock, &product.CreatedAt)

	utils.RespondSuccess(c, http.StatusOK, "Produk berhasil diupdate", product)
}

// DeleteProduct menghapus produk berdasarkan ID.
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	result, err := database.DB.Exec("DELETE FROM products WHERE product_id = ?", id)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Gagal menghapus produk: "+err.Error())
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		utils.RespondError(c, http.StatusNotFound, "Produk tidak ditemukan")
		return
	}

	utils.RespondSuccess(c, http.StatusOK, "Produk berhasil dihapus", nil)
}
