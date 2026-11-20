// ==================================================================
// models/user.go — Struct User untuk autentikasi
// Sistem Informasi Penjualan (Point of Sales)
// ==================================================================

package models

// User merepresentasikan entitas pengguna (admin/kasir).
type User struct {
	UserID       string `json:"user_id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"` // tidak di-expose ke JSON response
	FullName     string `json:"full_name"`
	Role         string `json:"role"` // "admin" | "kasir"
}

// LoginRequest merepresentasikan body request login.
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse merepresentasikan response setelah login berhasil.
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
