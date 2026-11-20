// ==================================================================
// handlers/auth_handler.go — Handler Autentikasi (Login)
// Sistem Informasi Penjualan (Point of Sales)
// ==================================================================
//
// Endpoint: POST /api/auth/login
// Request : { "username": "admin", "password": "admin123" }
// Response: { "status": "success", "data": { "token": "...", "user": {...} } }
// ==================================================================

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"pos-system/database"
	"pos-system/middleware"
	"pos-system/models"
	"pos-system/utils"
)

// Login menangani proses autentikasi user.
// Memverifikasi username dan password menggunakan bcrypt,
// lalu menghasilkan JWT token jika berhasil.
func Login(c *gin.Context) {
	var input models.LoginRequest

	// Bind JSON body
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Format request tidak valid")
		return
	}

	// Validasi input
	if validationErrors := utils.ValidateLoginInput(input); len(validationErrors) > 0 {
		utils.RespondValidationError(c, validationErrors)
		return
	}

	// Cari user di database berdasarkan username
	var user models.User
	err := database.DB.QueryRow(
		"SELECT user_id, username, password_hash, full_name, role FROM users WHERE username = ?",
		input.Username,
	).Scan(&user.UserID, &user.Username, &user.PasswordHash, &user.FullName, &user.Role)

	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, "Username atau password salah")
		return
	}

	// Verifikasi password menggunakan bcrypt
	if !utils.CheckPasswordHash(input.Password, user.PasswordHash) {
		utils.RespondError(c, http.StatusUnauthorized, "Username atau password salah")
		return
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(user)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Gagal membuat token")
		return
	}

	// Response sukses
	utils.RespondSuccess(c, http.StatusOK, "Login berhasil", models.LoginResponse{
		Token: token,
		User:  user,
	})
}
