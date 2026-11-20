// ==================================================================
// models/product.go — Struct Produk
// Sistem Informasi Penjualan (Point of Sales)
// ==================================================================

package models

// Product merepresentasikan entitas produk.
type Product struct {
	ProductID string  `json:"product_id"`
	Name      string  `json:"name"`
	Category  string  `json:"category"`
	Price     float64 `json:"price"`
	Stock     int     `json:"stock"`
	CreatedAt string  `json:"created_at"`
}

// ProductInput merepresentasikan body request untuk create/update produk.
type ProductInput struct {
	Name     string  `json:"name" binding:"required"`
	Category string  `json:"category" binding:"required"`
	Price    float64 `json:"price" binding:"required"`
	Stock    int     `json:"stock"`
}
