// ==================================================================
// models/transaction.go — Struct Transaksi
// Sistem Informasi Penjualan (Point of Sales)
// ==================================================================

package models

// Transaction merepresentasikan entitas transaksi penjualan.
type Transaction struct {
	TransactionID   string  `json:"transaction_id"`
	ProductID       string  `json:"product_id"`
	ProductName     string  `json:"product_name"`
	Quantity        int     `json:"quantity"`
	UnitPrice       float64 `json:"unit_price"`
	DiscountPercent float64 `json:"discount_percent"`
	TotalPrice      float64 `json:"total_price"`
	PaymentMethod   string  `json:"payment_method"`
	TransactionDate string  `json:"transaction_date"`
	CashierName     string  `json:"cashier_name"`
}

// TransactionInput merepresentasikan body request untuk create transaksi.
type TransactionInput struct {
	ProductID       string  `json:"product_id" binding:"required"`
	Quantity        int     `json:"quantity" binding:"required"`
	DiscountPercent float64 `json:"discount_percent"`
	PaymentMethod   string  `json:"payment_method" binding:"required"`
	CashierName     string  `json:"cashier_name" binding:"required"`
}

// APIResponse merupakan format response API yang seragam.
type APIResponse struct {
	Status  string        `json:"status"`            // "success" | "error"
	Message string        `json:"message,omitempty"`
	Data    interface{}   `json:"data,omitempty"`
	Errors  []FieldError  `json:"errors,omitempty"`
}

// FieldError merepresentasikan error validasi per-field.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
