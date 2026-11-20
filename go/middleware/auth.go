// ==================================================================
// middleware/auth.go — JWT Authentication Middleware
// Sistem Informasi Penjualan (Point of Sales)
// ==================================================================

package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"pos-system/models"
)

// JWTSecret adalah kunci rahasia untuk signing JWT token.
// Pada produksi, gunakan environment variable.
var JWTSecret = []byte("pos-system-secret-key-2026")

// Claims merupakan custom JWT claims yang menyimpan data user.
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken membuat JWT token baru untuk user yang berhasil login.
// Token berlaku selama 24 jam.
func GenerateToken(user models.User) (string, error) {
	claims := Claims{
		UserID:   user.UserID,
		Username: user.Username,
		FullName: user.FullName,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "pos-system",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}

// AuthMiddleware memvalidasi JWT token pada setiap request yang terproteksi.
// Token diambil dari header Authorization: Bearer <token>.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Status:  "error",
				Message: "Token tidak ditemukan. Silakan login terlebih dahulu.",
			})
			c.Abort()
			return
		}

		// Ambil token dari header "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Status:  "error",
				Message: "Format token tidak valid. Gunakan: Bearer <token>",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Parse dan validasi token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return JWTSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Status:  "error",
				Message: "Token tidak valid atau sudah expired. Silakan login ulang.",
			})
			c.Abort()
			return
		}

		// Simpan data user dari token ke context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("full_name", claims.FullName)
		c.Set("role", claims.Role)

		c.Next()
	}
}
