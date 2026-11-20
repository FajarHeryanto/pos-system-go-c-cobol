// ==================================================================
// utils/hash.go — Bcrypt Password Hashing Utilities
// Sistem Informasi Penjualan (Point of Sales)
// ==================================================================

package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword menghasilkan bcrypt hash dari password plaintext.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash memverifikasi password plaintext terhadap bcrypt hash.
// Mengembalikan true jika password cocok.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
