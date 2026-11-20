// ==================================================================
// utils/validator.go — Validasi Input
// Sistem Informasi Penjualan (Point of Sales)
// ==================================================================
//
// Validasi ganda (dual validation): dilakukan di frontend (JS) DAN
// di backend (Go) sebagai keamanan utama (Bab 9.3 RPD).
// ==================================================================

package utils

import (
	"regexp"
	"strings"

	"pos-system/models"
)

// ValidateProductInput memvalidasi input produk sesuai Bab 5.1.5 RPD.
func ValidateProductInput(input models.ProductInput) []models.FieldError {
	var errors []models.FieldError

	// Nama produk: non-empty, 2-50 karakter
	name := strings.TrimSpace(input.Name)
	if len(name) < 2 || len(name) > 50 {
		errors = append(errors, models.FieldError{
			Field:   "name",
			Message: "Nama produk harus 2-50 karakter",
		})
	}

	// Kategori: non-empty
	category := strings.TrimSpace(input.Category)
	if category == "" {
		errors = append(errors, models.FieldError{
			Field:   "category",
			Message: "Kategori tidak boleh kosong",
		})
	}

	// Harga: > 0, max 99.999.999
	if input.Price <= 0 {
		errors = append(errors, models.FieldError{
			Field:   "price",
			Message: "Harga harus > 0",
		})
	} else if input.Price > 99999999 {
		errors = append(errors, models.FieldError{
			Field:   "price",
			Message: "Harga maksimal Rp 99.999.999",
		})
	}

	// Stok: >= 0
	if input.Stock < 0 {
		errors = append(errors, models.FieldError{
			Field:   "stock",
			Message: "Stok tidak boleh negatif",
		})
	}

	return errors
}

// ValidateTransactionInput memvalidasi input transaksi sesuai Bab 5.1.5 RPD.
func ValidateTransactionInput(input models.TransactionInput) []models.FieldError {
	var errors []models.FieldError

	// Product ID: non-empty
	if strings.TrimSpace(input.ProductID) == "" {
		errors = append(errors, models.FieldError{
			Field:   "product_id",
			Message: "Product ID tidak boleh kosong",
		})
	}

	// Quantity: > 0
	if input.Quantity <= 0 {
		errors = append(errors, models.FieldError{
			Field:   "quantity",
			Message: "Quantity harus > 0",
		})
	}

	// Diskon: 0-100
	if input.DiscountPercent < 0 || input.DiscountPercent > 100 {
		errors = append(errors, models.FieldError{
			Field:   "discount_percent",
			Message: "Diskon harus antara 0-100",
		})
	}

	// Payment method: hanya Cash, Debit, Credit
	validMethods := map[string]bool{"Cash": true, "Debit": true, "Credit": true}
	if !validMethods[input.PaymentMethod] {
		errors = append(errors, models.FieldError{
			Field:   "payment_method",
			Message: "Metode bayar harus Cash, Debit, atau Credit",
		})
	}

	// Cashier name: non-empty
	if strings.TrimSpace(input.CashierName) == "" {
		errors = append(errors, models.FieldError{
			Field:   "cashier_name",
			Message: "Nama kasir tidak boleh kosong",
		})
	}

	return errors
}

// ValidateLoginInput memvalidasi input login sesuai Bab 5.1.5 RPD.
func ValidateLoginInput(input models.LoginRequest) []models.FieldError {
	var errors []models.FieldError

	// Username: 3-20 karakter, alfanumerik
	username := strings.TrimSpace(input.Username)
	if len(username) < 3 || len(username) > 20 {
		errors = append(errors, models.FieldError{
			Field:   "username",
			Message: "Username harus 3-20 karakter",
		})
	} else {
		matched, _ := regexp.MatchString(`^[a-zA-Z0-9]+$`, username)
		if !matched {
			errors = append(errors, models.FieldError{
				Field:   "username",
				Message: "Username hanya boleh alfanumerik",
			})
		}
	}

	// Password: min 6 karakter
	if len(input.Password) < 6 {
		errors = append(errors, models.FieldError{
			Field:   "password",
			Message: "Password minimal 6 karakter",
		})
	}

	return errors
}
